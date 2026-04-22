// SPDX-License-Identifier: MIT
#ifndef GO_QUICKJS_BRIDGE_H
#define GO_QUICKJS_BRIDGE_H

#include "../../third_party/quickjs/quickjs.h"
#include <stddef.h>
#include <stdint.h>

JSRuntime *qjs_new_runtime(void);
void qjs_free_runtime(JSRuntime *rt);
void qjs_set_memory_limit(JSRuntime *rt, size_t limit_bytes);

JSContext *qjs_new_context(JSRuntime *rt);
void qjs_free_context(JSContext *ctx);

JSValue qjs_eval(JSContext *ctx, const char *code, size_t code_len, const char *filename, int eval_flags);
int qjs_is_exception(JSValueConst v);
JSValue qjs_get_exception(JSContext *ctx);

JSValue qjs_dup_value(JSContext *ctx, JSValueConst v);
void qjs_free_value(JSContext *ctx, JSValue v);

const char *qjs_to_cstring(JSContext *ctx, JSValueConst v);
void qjs_free_cstring(JSContext *ctx, const char *str);

int qjs_to_int32(JSContext *ctx, int32_t *out, JSValueConst v);
int qjs_to_float64(JSContext *ctx, double *out, JSValueConst v);
int qjs_to_bool(JSContext *ctx, JSValueConst v);

JSValue qjs_new_int32(JSContext *ctx, int32_t v);
JSValue qjs_new_float64(JSContext *ctx, double v);
JSValue qjs_new_bool(JSContext *ctx, int v);
JSValue qjs_new_string(JSContext *ctx, const char *str);
JSValue qjs_new_null(void);
JSValue qjs_new_undefined(void);
JSValue qjs_new_object(JSContext *ctx);
JSValue qjs_new_array(JSContext *ctx);

int qjs_is_null(JSValueConst v);
int qjs_is_undefined(JSValueConst v);
int qjs_is_bool(JSValueConst v);
int qjs_is_number(JSValueConst v);
int qjs_is_string(JSValueConst v);
int qjs_is_object(JSValueConst v);
int qjs_is_array(JSContext *ctx, JSValueConst v);
int qjs_is_function(JSContext *ctx, JSValueConst v);

JSValue qjs_get_global_object(JSContext *ctx);
int qjs_set_global_function(JSContext *ctx, const char *name, int length, int magic);

JSValue qjs_get_property_str(JSContext *ctx, JSValueConst obj, const char *prop);
int qjs_set_property_str_dup(JSContext *ctx, JSValue obj, const char *prop, JSValueConst value);
JSValue qjs_get_property_uint32(JSContext *ctx, JSValueConst obj, uint32_t idx);
int qjs_set_property_uint32_dup(JSContext *ctx, JSValue obj, uint32_t idx, JSValueConst value);

JSValue qjs_throw_type_error(JSContext *ctx, const char *message);

#endif
