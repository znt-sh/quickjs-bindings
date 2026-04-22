// SPDX-License-Identifier: MIT
#include "bridge.h"

extern JSValue goQuickjsInvoke(JSContext *ctx, JSValueConst this_val, int argc, JSValueConst *argv, int magic);

static JSValue qjs_go_function(JSContext *ctx, JSValueConst this_val, int argc, JSValueConst *argv, int magic) {
	return goQuickjsInvoke(ctx, this_val, argc, argv, magic);
}

JSRuntime *qjs_new_runtime(void) {
	return JS_NewRuntime();
}

void qjs_free_runtime(JSRuntime *rt) {
	JS_FreeRuntime(rt);
}

void qjs_set_memory_limit(JSRuntime *rt, size_t limit_bytes) {
	JS_SetMemoryLimit(rt, limit_bytes);
}

JSContext *qjs_new_context(JSRuntime *rt) {
	return JS_NewContext(rt);
}

void qjs_free_context(JSContext *ctx) {
	JS_FreeContext(ctx);
}

JSValue qjs_eval(JSContext *ctx, const char *code, size_t code_len, const char *filename, int eval_flags) {
	return JS_Eval(ctx, code, code_len, filename, eval_flags);
}

int qjs_is_exception(JSValueConst v) {
	return JS_IsException(v);
}

JSValue qjs_get_exception(JSContext *ctx) {
	return JS_GetException(ctx);
}

JSValue qjs_dup_value(JSContext *ctx, JSValueConst v) {
	return JS_DupValue(ctx, v);
}

void qjs_free_value(JSContext *ctx, JSValue v) {
	JS_FreeValue(ctx, v);
}

const char *qjs_to_cstring(JSContext *ctx, JSValueConst v) {
	return JS_ToCString(ctx, v);
}

void qjs_free_cstring(JSContext *ctx, const char *str) {
	JS_FreeCString(ctx, str);
}

int qjs_to_int32(JSContext *ctx, int32_t *out, JSValueConst v) {
	return JS_ToInt32(ctx, out, v);
}

int qjs_to_float64(JSContext *ctx, double *out, JSValueConst v) {
	return JS_ToFloat64(ctx, out, v);
}

int qjs_to_bool(JSContext *ctx, JSValueConst v) {
	return JS_ToBool(ctx, v);
}

JSValue qjs_new_int32(JSContext *ctx, int32_t v) {
	return JS_NewInt32(ctx, v);
}

JSValue qjs_new_float64(JSContext *ctx, double v) {
	return JS_NewFloat64(ctx, v);
}

JSValue qjs_new_bool(JSContext *ctx, int v) {
	return JS_NewBool(ctx, v);
}

JSValue qjs_new_string(JSContext *ctx, const char *str) {
	return JS_NewString(ctx, str);
}

JSValue qjs_new_null(void) {
	return JS_NULL;
}

JSValue qjs_new_undefined(void) {
	return JS_UNDEFINED;
}

JSValue qjs_new_object(JSContext *ctx) {
	return JS_NewObject(ctx);
}

JSValue qjs_new_array(JSContext *ctx) {
	return JS_NewArray(ctx);
}

int qjs_is_null(JSValueConst v) {
	return JS_IsNull(v);
}

int qjs_is_undefined(JSValueConst v) {
	return JS_IsUndefined(v);
}

int qjs_is_bool(JSValueConst v) {
	return JS_IsBool(v);
}

int qjs_is_number(JSValueConst v) {
	return JS_IsNumber(v);
}

int qjs_is_string(JSValueConst v) {
	return JS_IsString(v);
}

int qjs_is_object(JSValueConst v) {
	return JS_IsObject(v);
}

int qjs_is_array(JSContext *ctx, JSValueConst v) {
	return JS_IsArray(ctx, v);
}

int qjs_is_function(JSContext *ctx, JSValueConst v) {
	return JS_IsFunction(ctx, v);
}

JSValue qjs_get_global_object(JSContext *ctx) {
	return JS_GetGlobalObject(ctx);
}

int qjs_set_global_function(JSContext *ctx, const char *name, int length, int magic) {
	JSValue global = JS_GetGlobalObject(ctx);
	JSValue fn = JS_NewCFunctionMagic(ctx, qjs_go_function, name, length, JS_CFUNC_generic_magic, magic);
	int rc = JS_SetPropertyStr(ctx, global, name, fn);
	JS_FreeValue(ctx, global);
	return rc;
}

JSValue qjs_get_property_str(JSContext *ctx, JSValueConst obj, const char *prop) {
	return JS_GetPropertyStr(ctx, obj, prop);
}

int qjs_set_property_str_dup(JSContext *ctx, JSValue obj, const char *prop, JSValueConst value) {
	return JS_SetPropertyStr(ctx, obj, prop, JS_DupValue(ctx, value));
}

JSValue qjs_get_property_uint32(JSContext *ctx, JSValueConst obj, uint32_t idx) {
	return JS_GetPropertyUint32(ctx, obj, idx);
}

int qjs_set_property_uint32_dup(JSContext *ctx, JSValue obj, uint32_t idx, JSValueConst value) {
	return JS_SetPropertyUint32(ctx, obj, idx, JS_DupValue(ctx, value));
}

JSValue qjs_throw_type_error(JSContext *ctx, const char *message) {
	return JS_ThrowTypeError(ctx, "%s", message);
}
