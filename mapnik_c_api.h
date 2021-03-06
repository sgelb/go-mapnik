#include "stdint.h"
#ifndef MAPNIK_C_API_H
#define MAPNIK_C_API_H

#if defined(WIN32) || defined(WINDOWS) || defined(_WIN32) || defined(_WINDOWS)
#  define MAPNIKCAPICALL __declspec(dllexport)
#else
#  define MAPNIKCAPICALL
#endif

#ifdef __cplusplus
extern "C"
{
#endif

extern const int mapnik_version;
extern const char *mapnik_version_string;
extern const int mapnik_version_major;
extern const int mapnik_version_minor;
extern const int mapnik_version_patch;

MAPNIKCAPICALL int mapnik_register_datasources(const char* path);
MAPNIKCAPICALL int mapnik_register_fonts(const char* path);

const int MAPNIK_NONE = 0;
const int MAPNIK_DEBUG = 1;
const int MAPNIK_WARN = 2;
const int MAPNIK_ERROR = 3;

MAPNIKCAPICALL void mapnik_logging_set_severity(int);

MAPNIKCAPICALL const char * mapnik_register_last_error();

// BBOX
typedef struct _mapnik_bbox_t mapnik_bbox_t;
MAPNIKCAPICALL mapnik_bbox_t * mapnik_bbox(double minx, double miny, double maxx, double maxy);
MAPNIKCAPICALL void mapnik_bbox_free(mapnik_bbox_t * b);


// Image
MAPNIKCAPICALL typedef struct _mapnik_image_t mapnik_image_t;
MAPNIKCAPICALL void mapnik_image_free(mapnik_image_t * i);

MAPNIKCAPICALL const char * mapnik_image_last_error(mapnik_image_t * i);

typedef struct _mapnik_image_blob_t {
    char *ptr;
    unsigned int len;
} mapnik_image_blob_t;
MAPNIKCAPICALL void mapnik_image_blob_free(mapnik_image_blob_t * b);

MAPNIKCAPICALL mapnik_image_blob_t * mapnik_image_to_blob(mapnik_image_t * i, const char * format);

MAPNIKCAPICALL const uint8_t * mapnik_image_to_raw(mapnik_image_t * i, size_t *size);
MAPNIKCAPICALL mapnik_image_t * mapnik_image_from_raw(const uint8_t * raw, int width, int height);


// Parameters
typedef struct _mapnik_parameters_t mapnik_parameters_t;

MAPNIKCAPICALL mapnik_parameters_t *mapnik_parameters();

MAPNIKCAPICALL void mapnik_parameters_free(mapnik_parameters_t *p);

MAPNIKCAPICALL void mapnik_parameters_set(mapnik_parameters_t *p, const char *key, const char *value);


// Datasource
typedef struct _mapnik_datasource_t mapnik_datasource_t;

MAPNIKCAPICALL mapnik_datasource_t *mapnik_datasource(mapnik_parameters_t *p);

MAPNIKCAPICALL void mapnik_datasource_free(mapnik_datasource_t *ds);


// Layer
typedef struct _mapnik_layer_t mapnik_layer_t;

MAPNIKCAPICALL mapnik_layer_t *mapnik_layer(const char *name, const char *srs);
MAPNIKCAPICALL void mapnik_layer_free(mapnik_layer_t *l);

MAPNIKCAPICALL void mapnik_layer_add_style(mapnik_layer_t *l, const char *stylename);

MAPNIKCAPICALL void mapnik_layer_set_datasource(mapnik_layer_t *l, mapnik_datasource_t *ds);


//  Map
typedef struct _mapnik_map_t mapnik_map_t;

MAPNIKCAPICALL mapnik_map_t * mapnik_map(unsigned int width, unsigned int height);
MAPNIKCAPICALL void mapnik_map_free(mapnik_map_t * m);

MAPNIKCAPICALL const char * mapnik_map_last_error(mapnik_map_t * m);

MAPNIKCAPICALL int mapnik_map_load(mapnik_map_t * m, const char* stylesheet);
MAPNIKCAPICALL int mapnik_map_load_string(mapnik_map_t *m, const char* s, const char* base_path);

MAPNIKCAPICALL const char * mapnik_map_get_srs(mapnik_map_t * m);
MAPNIKCAPICALL int mapnik_map_set_srs(mapnik_map_t * m, const char* srs);
MAPNIKCAPICALL int mapnik_map_set_aspect_fix_mode(mapnik_map_t * m, int afm);
MAPNIKCAPICALL int mapnik_map_get_aspect_fix_mode(mapnik_map_t * m);

MAPNIKCAPICALL void mapnik_map_resize(mapnik_map_t * m, unsigned int width, unsigned int height);
MAPNIKCAPICALL double mapnik_map_get_scale_denominator(mapnik_map_t * m);
MAPNIKCAPICALL void mapnik_map_set_buffer_size(mapnik_map_t * m, int buffer_size);

MAPNIKCAPICALL int mapnik_map_background(mapnik_map_t * m, uint8_t *r, uint8_t *g, uint8_t *b, uint8_t *a);
MAPNIKCAPICALL void mapnik_map_set_background(mapnik_map_t * m, uint8_t r, uint8_t g, uint8_t b, uint8_t a);

MAPNIKCAPICALL int mapnik_map_zoom_all(mapnik_map_t * m);
MAPNIKCAPICALL void mapnik_map_zoom_to_box(mapnik_map_t * m, mapnik_bbox_t * b);

MAPNIKCAPICALL void mapnik_map_set_maximum_extent(mapnik_map_t * m, double x0, double y0, double x1, double y1);
MAPNIKCAPICALL void mapnik_map_reset_maximum_extent(mapnik_map_t * m);

MAPNIKCAPICALL int mapnik_map_render_to_file(mapnik_map_t * m, const char* filepath, double scale, double scale_factor, const char *format);
MAPNIKCAPICALL mapnik_image_t * mapnik_map_render_to_image(mapnik_map_t * m, double scale, double scale_factor);

MAPNIKCAPICALL void mapnik_map_add_layer(mapnik_map_t *m, mapnik_layer_t *l);

MAPNIKCAPICALL int mapnik_map_layer_count(mapnik_map_t * m);
MAPNIKCAPICALL const char * mapnik_map_layer_name(mapnik_map_t * m, size_t idx);
MAPNIKCAPICALL int mapnik_map_layer_is_active(mapnik_map_t * m, size_t idx);
MAPNIKCAPICALL void mapnik_map_layer_set_active(mapnik_map_t * m, size_t idx, int active);

#ifdef __cplusplus
}
#endif

#endif // MAPNIK_C_API_H
