#ifndef __cv_handler_h__
#define __cv_handler_h__

#include <opencv/highgui.h>
#include <opencv/cv.h>
#include <stdio.h>
#include <strings.h>
#include <stdlib.h>
#include <stdint.h>

typedef struct PixelDim{
    int64_t width;
    int64_t height;
} PixelDim;

typedef struct Blob{
    unsigned char *data;
    unsigned int length;    
} Blob;

Blob *resizer(unsigned char *data, Blob *in, PixelDim *zoom, int quality, int method, const char *format, CvRect *roi);
Blob *blender(const Blob *bg, const Blob *fg, const Blob *mask, int quality, const char *format, const float alpha, CvRect *roi);

#endif
