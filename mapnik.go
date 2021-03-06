// Package mapnik renders beautiful maps with Mapnik.
package mapnik

//go:generate bash ./configure.bash

// #include <stdlib.h>
// #include "mapnik_c_api.h"
import "C"

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"unsafe"
)

type LogLevel int

var (
	None  = LogLevel(C.MAPNIK_NONE)
	Debug = LogLevel(C.MAPNIK_DEBUG)
	Warn  = LogLevel(C.MAPNIK_WARN)
	Error = LogLevel(C.MAPNIK_ERROR)
)

func init() {
	// register default datasources path and fonts path like the python bindings do
	RegisterDatasources(pluginPath)
	RegisterFonts(fontPath)
}

// RegisterDatasources adds path to the Mapnik plugin search path.
func RegisterDatasources(path string) error {
	cs := C.CString(path)
	defer C.free(unsafe.Pointer(cs))
	if C.mapnik_register_datasources(cs) != 0 {
		return errors.New("mapnik: " + C.GoString(C.mapnik_register_last_error()))
	}
	return nil
}

// RegisterFonts adds path to the Mapnik fonts search path.
func RegisterFonts(path string) error {
	cs := C.CString(path)
	defer C.free(unsafe.Pointer(cs))
	if C.mapnik_register_fonts(cs) != 0 {
		return errors.New("mapnik: " + C.GoString(C.mapnik_register_last_error()))
	}
	return nil
}

// LogSeverity sets the global log level for Mapnik. Requires a Mapnik build with logging enabled.
func LogSeverity(level LogLevel) {
	C.mapnik_logging_set_severity(C.int(level))
}

type version struct {
	Numeric int
	Major   int
	Minor   int
	Patch   int
	String  string
}

var Version version

func init() {
	Version.Numeric = int(C.mapnik_version)
	Version.Major = int(C.mapnik_version_major)
	Version.Minor = int(C.mapnik_version_minor)
	Version.Patch = int(C.mapnik_version_patch)
	Version.String = C.GoString(C.mapnik_version_string)
}

// Datasource base type
type Datasource struct {
	ds *C.struct__mapnik_datasource_t
}

// NewDatasource initializes a new Datasource
func NewDatasource(params map[string]string) *Datasource {
	p := C.mapnik_parameters()
	defer C.mapnik_parameters_free(p)
	for k, v := range params {
		kcs := C.CString(k)
		defer C.free(unsafe.Pointer(kcs))
		vcs := C.CString(v)
		defer C.free(unsafe.Pointer(vcs))
		C.mapnik_parameters_set(p, kcs, vcs)
	}
	return &Datasource{C.mapnik_datasource(p)}
}

// Free deallocates the datasource.
func (ds *Datasource) Free() {
	C.mapnik_datasource_free(ds.ds)
	ds.ds = nil
}

// Layer base type
type Layer struct {
	l *C.struct__mapnik_layer_t
}

// NewLayer initializes a new Layer
func NewLayer(name string, srs string) *Layer {
	namecs := C.CString(name)
	defer C.free(unsafe.Pointer(namecs))
	srscs := C.CString(srs)
	defer C.free(unsafe.Pointer(srscs))
	return &Layer{C.mapnik_layer(namecs, srscs)}
}

// Free deallocates the layer.
func (l *Layer) Free() {
	C.mapnik_layer_free(l.l)
	l.l = nil
}

// AddStyle adds a style.
func (l *Layer) AddStyle(stylename string) {
	cs := C.CString(stylename)
	defer C.free(unsafe.Pointer(cs))
	C.mapnik_layer_add_style(l.l, cs)
}

// SetDatasource sets the datasource.
func (l *Layer) SetDatasource(ds *Datasource) {
	C.mapnik_layer_set_datasource(l.l, ds.ds)
}

// Map base type
type Map struct {
	m           *C.struct__mapnik_map_t
	width       int
	height      int
	layerStatus []bool
}

// New initializes a new Map.
func New() *Map {
	return &Map{
		m:      C.mapnik_map(C.uint(800), C.uint(600)),
		width:  800,
		height: 600,
	}
}

func (m *Map) lastError() error {
	return errors.New("mapnik: " + C.GoString(C.mapnik_map_last_error(m.m)))
}

// Load reads in a Mapnik map XML.
func (m *Map) Load(stylesheet string) error {
	cs := C.CString(stylesheet)
	defer C.free(unsafe.Pointer(cs))
	if C.mapnik_map_load(m.m, cs) != 0 {
		return m.lastError()
	}
	return nil
}

// LoadString reads in a Mapnik map from a XML string.
func (m *Map) LoadString(s string, basePath string) error {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	bs := C.CString(basePath)
	defer C.free(unsafe.Pointer(bs))
	if C.mapnik_map_load_string(m.m, cs, bs) != 0 {
		return m.lastError()
	}
	return nil
}

// Resize changes the map size in pixel.
func (m *Map) Resize(width, height int) {
	C.mapnik_map_resize(m.m, C.uint(width), C.uint(height))
	m.width = width
	m.height = height
}

// Free deallocates the map.
func (m *Map) Free() {
	C.mapnik_map_free(m.m)
	m.m = nil
}

// SRS returns the projection of the map.
func (m *Map) SRS() string {
	return C.GoString(C.mapnik_map_get_srs(m.m))
}

// SetSRS sets the projection of the map as a proj4 string ('+init=epsg:4326', etc).
func (m *Map) SetSRS(srs string) {
	cs := C.CString(srs)
	defer C.free(unsafe.Pointer(cs))
	C.mapnik_map_set_srs(m.m, cs)
}

// FixMode defines what mapnik does when the aspect ratio of BBox and map are not 1:1.
type FixMode int

const (
	// GrowBBox grows the width or height of the specified geo bbox to fill the map size. Default behaviour.
	GrowBBox FixMode = iota
	// GrowCanvas grows the width or height of the map to accomodate the specified geo bbox.
	GrowCanvas
	// ShrinkBBox shrinks the width or height of the specified geo bbox to fill the map size.
	ShrinkBBox
	// ShrinkCanvas shrinks the width or height of the map to accomodate the specified geo bbox.
	ShrinkCanvas
	// AdjustBBoxWidth adjusts the width of the specified geo bbox, leaves height and map size unchanged.
	AdjustBBoxWidth
	// AdjustBBoxHeight adjusts the height of the specified geo bbox, leaves width and map size unchanged.
	AdjustBBoxHeight
	// AdjustCanvasWidth adjusts the width of the map, leaves height and geo bbox unchanged.
	AdjustCanvasWidth
	// AdjustCanvasHeight adjusts the height of the map, leaves width and geo bbox unchanged.
	AdjustCanvasHeight
	// Respect does nothing
	Respect
)

// SetAspectFixMode sets the aspect fix mode. Set before Resize and ZoomAll/ZoomTo.
func (m *Map) SetAspectFixMode(f FixMode) error {
	if C.mapnik_map_set_aspect_fix_mode(m.m, C.int(f)) != 0 {
		return errors.New("mapnik: " + C.GoString(C.mapnik_register_last_error()))
	}
	return nil
}

// AspectFixMode returns the current aspect fix mode.
func (m *Map) AspectFixMode() FixMode {
	return FixMode(C.mapnik_map_get_aspect_fix_mode(m.m))
}

// ScaleDenominator returns the current scale denominator. Call after Resize and ZoomAll/ZoomTo.
func (m *Map) ScaleDenominator() float64 {
	return float64(C.mapnik_map_get_scale_denominator(m.m))
}

// ZoomAll zooms to the maximum extent.
func (m *Map) ZoomAll() error {
	if C.mapnik_map_zoom_all(m.m) != 0 {
		return m.lastError()
	}
	return nil
}

// ZoomTo zooms to the given bounding box.
func (m *Map) ZoomTo(minx, miny, maxx, maxy float64) {
	bbox := C.mapnik_bbox(C.double(minx), C.double(miny), C.double(maxx), C.double(maxy))
	defer C.mapnik_bbox_free(bbox)
	C.mapnik_map_zoom_to_box(m.m, bbox)
}

func (m *Map) BackgroundColor() color.NRGBA {
	c := color.NRGBA{}
	C.mapnik_map_background(m.m, (*C.uint8_t)(&c.R), (*C.uint8_t)(&c.G), (*C.uint8_t)(&c.B), (*C.uint8_t)(&c.A))
	return c
}

func (m *Map) SetBackgroundColor(c color.NRGBA) {
	C.mapnik_map_set_background(m.m, C.uint8_t(c.R), C.uint8_t(c.G), C.uint8_t(c.B), C.uint8_t(c.A))
}

func (m *Map) printLayerStatus() {
	n := m.CountLayers()
	for i := 0; i < n; i++ {
		fmt.Println(
			C.GoString(C.mapnik_map_layer_name(m.m, C.size_t(i))),
			C.mapnik_map_layer_is_active(m.m, C.size_t(i)),
		)
	}
}

func (m *Map) storeLayerStatus() {
	if len(m.layerStatus) > 0 {
		return // allready stored
	}
	m.layerStatus = m.currentLayerStatus()
}

func (m *Map) currentLayerStatus() []bool {
	n := m.CountLayers()
	active := make([]bool, n)
	for i := 0; i < n; i++ {
		if C.mapnik_map_layer_is_active(m.m, C.size_t(i)) == 1 {
			active[i] = true
		}
	}
	return active
}

func (m *Map) resetLayerStatus() {
	if len(m.layerStatus) == 0 {
		return // not stored
	}
	n := m.CountLayers()
	if n > len(m.layerStatus) {
		// should not happen
		return
	}
	for i := 0; i < n; i++ {
		if m.layerStatus[i] {
			C.mapnik_map_layer_set_active(m.m, C.size_t(i), 1)
		} else {
			C.mapnik_map_layer_set_active(m.m, C.size_t(i), 0)
		}
	}
	m.layerStatus = nil
}

// Status defines if a layer should be rendered or not.
type Status int

const (
	// Exclude layer from rendering.
	Exclude Status = -1
	// Default keeps layer at the current rendering status.
	Default Status = 0
	// Include layer for rendering.
	Include Status = 1
)

type LayerSelector interface {
	Select(layername string) Status
}

type SelectorFunc func(string) Status

func (f SelectorFunc) Select(layername string) Status {
	return f(layername)
}

// AddLayer adds a layer.
func (m *Map) AddLayer(l *Layer) {
	C.mapnik_map_add_layer(m.m, l.l)
}

// SelectLayers enables/disables single layers. LayerSelector or SelectorFunc gets called for each layer.
func (m *Map) SelectLayers(selector LayerSelector) {
	m.storeLayerStatus()
	n := m.CountLayers()
	for i := 0; i < n; i++ {
		layerName := C.GoString(C.mapnik_map_layer_name(m.m, C.size_t(i)))
		switch selector.Select(layerName) {
		case Include:
			C.mapnik_map_layer_set_active(m.m, C.size_t(i), 1)
		case Exclude:
			C.mapnik_map_layer_set_active(m.m, C.size_t(i), 0)
		case Default:
		}
	}
}

// CountLayers returns count of layers
func (m *Map) CountLayers() int {
	return int(C.mapnik_map_layer_count(m.m))
}

// ResetLayers resets all layers to the initial status.
func (m *Map) ResetLayers() {
	m.resetLayerStatus()
}

func (m *Map) SetMaxExtent(minx, miny, maxx, maxy float64) {
	C.mapnik_map_set_maximum_extent(m.m, C.double(minx), C.double(miny), C.double(maxx), C.double(maxy))
}

func (m *Map) ResetMaxExtent() {
	C.mapnik_map_reset_maximum_extent(m.m)
}

// RenderOpts defines rendering options.
type RenderOpts struct {
	// Scale renders the map at a fixed scale denominator.
	Scale float64
	// ScaleFactor renders the map with larger fonts sizes, line width, etc. For printing or retina/hq iamges.
	ScaleFactor float64
	// Format for the rendered image ('jpeg80', 'png256', etc. see: https://github.com/mapnik/mapnik/wiki/Image-IO)
	Format string
}

// Render returns the map as an encoded image.
func (m *Map) Render(opts RenderOpts) ([]byte, error) {
	scaleFactor := opts.ScaleFactor
	if scaleFactor == 0.0 {
		scaleFactor = 1.0
	}
	i := C.mapnik_map_render_to_image(m.m, C.double(opts.Scale), C.double(scaleFactor))
	if i == nil {
		return nil, m.lastError()
	}
	defer C.mapnik_image_free(i)
	if opts.Format == "raw" {
		size := 0
		raw := C.mapnik_image_to_raw(i, (*C.size_t)(unsafe.Pointer(&size)))
		return C.GoBytes(unsafe.Pointer(raw), C.int(size)), nil
	}
	var format *C.char
	if opts.Format != "" {
		format = C.CString(opts.Format)
	} else {
		format = C.CString("png256")
	}
	b := C.mapnik_image_to_blob(i, format)
	if b == nil {
		return nil, errors.New("mapnik: " + C.GoString(C.mapnik_image_last_error(i)))
	}
	C.free(unsafe.Pointer(format))
	defer C.mapnik_image_blob_free(b)
	return C.GoBytes(unsafe.Pointer(b.ptr), C.int(b.len)), nil
}

// RenderImage returns the map as an unencoded image.Image.
func (m *Map) RenderImage(opts RenderOpts) (*image.NRGBA, error) {
	scaleFactor := opts.ScaleFactor
	if scaleFactor == 0.0 {
		scaleFactor = 1.0
	}
	i := C.mapnik_map_render_to_image(m.m, C.double(opts.Scale), C.double(scaleFactor))
	if i == nil {
		return nil, m.lastError()
	}
	defer C.mapnik_image_free(i)
	size := 0
	raw := C.mapnik_image_to_raw(i, (*C.size_t)(unsafe.Pointer(&size)))
	b := C.GoBytes(unsafe.Pointer(raw), C.int(size))
	img := &image.NRGBA{
		Pix:    b,
		Stride: int(m.width * 4),
		Rect:   image.Rect(0, 0, int(m.width), int(m.height)),
	}
	return img, nil
}

// RenderToFile writes the map as an encoded image to the file system.
func (m *Map) RenderToFile(opts RenderOpts, path string) error {
	scaleFactor := opts.ScaleFactor
	if scaleFactor == 0.0 {
		scaleFactor = 1.0
	}
	cs := C.CString(path)
	defer C.free(unsafe.Pointer(cs))
	var format *C.char
	if opts.Format != "" {
		format = C.CString(opts.Format)
	} else {
		format = C.CString("png256")
	}
	defer C.free(unsafe.Pointer(format))
	if C.mapnik_map_render_to_file(m.m, cs, C.double(opts.Scale), C.double(scaleFactor), format) != 0 {
		return m.lastError()
	}
	return nil
}

// SetBufferSize sets the pixel buffer at the map image edges where Mapnik should not render any labels.
func (m *Map) SetBufferSize(s int) {
	C.mapnik_map_set_buffer_size(m.m, C.int(s))
}

// Encode image.Image with Mapniks image encoder.
func Encode(img image.Image, format string) ([]byte, error) {
	var i *C.mapnik_image_t
	switch img := img.(type) {
	// XXX does mapnik expect NRGBA or RGBA? this might stop working
	//as expected if we start encoding images with full alpha channel
	case *image.NRGBA:
		i = C.mapnik_image_from_raw(
			(*C.uint8_t)(unsafe.Pointer(&img.Pix[0])),
			C.int(img.Bounds().Dx()),
			C.int(img.Bounds().Dy()),
		)
	case *image.RGBA:
		i = C.mapnik_image_from_raw(
			(*C.uint8_t)(unsafe.Pointer(&img.Pix[0])),
			C.int(img.Bounds().Dx()),
			C.int(img.Bounds().Dy()),
		)
	}

	if i == nil {
		return nil, errors.New("unable to create image from raw")
	}
	defer C.mapnik_image_free(i)

	cformat := C.CString(format)
	b := C.mapnik_image_to_blob(i, cformat)
	if b == nil {
		return nil, errors.New("mapnik: " + C.GoString(C.mapnik_image_last_error(i)))
	}
	C.free(unsafe.Pointer(cformat))
	defer C.mapnik_image_blob_free(b)
	return C.GoBytes(unsafe.Pointer(b.ptr), C.int(b.len)), nil
}
