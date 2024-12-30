package zlog

import (
	"encoding/base64"
	"fmt"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"math"
	"sync"
	"time"
)

var (
	_bufferPool = buffer.NewPool()
	// 构建编码器 池子
	_keyValueEncoderPool = sync.Pool{New: func() interface{} {
		return &keyValueEncoder{}
	}}
)

// 自定义 Key=Value 的日志格式
type keyValueEncoder struct {
	*zapcore.EncoderConfig
	buf *buffer.Buffer

	// for encoding generic values by reflection
	reflectBuf *buffer.Buffer
	reflectEnc zapcore.ReflectedEncoder
}

func (enc *keyValueEncoder) AddArray(key string, arr zapcore.ArrayMarshaler) error {
	enc.addKey(key)
	return enc.AppendArray(arr)
}

func (enc *keyValueEncoder) AddObject(key string, obj zapcore.ObjectMarshaler) error {
	enc.addKey(key)
	return enc.AppendObject(obj)
}

func (enc *keyValueEncoder) AddBinary(key string, val []byte) {
	enc.AddString(key, base64.StdEncoding.EncodeToString(val))
}

func (enc *keyValueEncoder) AddByteString(key string, val []byte) {
	enc.addKey(key)
	enc.AppendByteString(val)
}

func (enc *keyValueEncoder) AddBool(key string, val bool) {
	enc.addKey(key)
	enc.AppendBool(val)
}

func (enc *keyValueEncoder) AddComplex128(key string, val complex128) {
	enc.addKey(key)
	enc.AppendComplex128(val)
}

func (enc *keyValueEncoder) AddComplex64(key string, val complex64) {
	enc.addKey(key)
	enc.AppendComplex64(val)
}

func (enc *keyValueEncoder) AddDuration(key string, val time.Duration) {
	enc.addKey(key)
	enc.AppendDuration(val)
}

func (enc *keyValueEncoder) AddFloat64(key string, val float64) {
	enc.addKey(key)
	enc.AppendFloat64(val)
}

func (enc *keyValueEncoder) AddFloat32(key string, val float32) {
	enc.addKey(key)
	enc.AppendFloat32(val)
}

func (enc *keyValueEncoder) AddInt64(key string, val int64) {
	enc.addKey(key)
	enc.AppendInt64(val)
}

func (enc *keyValueEncoder) resetReflectBuf() {
	if enc.reflectBuf == nil {
		enc.reflectBuf = _bufferPool.Get()
		enc.reflectEnc = enc.NewReflectedEncoder(enc.reflectBuf)
	} else {
		enc.reflectBuf.Reset()
	}
}

var nullLiteralBytes = []byte("null")

// Only invoke the standard JSON encoder if there is actually something to
// encode; otherwise write JSON null literal directly.
func (enc *keyValueEncoder) encodeReflected(obj interface{}) ([]byte, error) {
	if obj == nil {
		return nullLiteralBytes, nil
	}
	enc.resetReflectBuf()
	if err := enc.reflectEnc.Encode(obj); err != nil {
		return nil, err
	}
	enc.reflectBuf.TrimNewline()
	return enc.reflectBuf.Bytes(), nil
}

func (enc *keyValueEncoder) AddReflected(key string, obj interface{}) error {
	valueBytes, err := enc.encodeReflected(obj)
	if err != nil {
		return err
	}
	enc.addKey(key)
	_, err = enc.buf.Write(valueBytes)
	return err
}

func (enc *keyValueEncoder) OpenNamespace(key string) {

}

func (enc *keyValueEncoder) AddString(key, val string) {
	enc.addKey(key)
	enc.AppendString(val)
}

func (enc *keyValueEncoder) AddTime(key string, val time.Time) {
	enc.addKey(key)
	enc.AppendTime(val)
}

func (enc *keyValueEncoder) AddUint64(key string, val uint64) {
	enc.addKey(key)
	enc.AppendUint64(val)
}

func (enc *keyValueEncoder) AppendArray(arr zapcore.ArrayMarshaler) error {
	enc.addElementSeparator()
	enc.buf.AppendByte('[')
	err := arr.MarshalLogArray(enc)
	enc.buf.AppendByte(']')
	return err
}

func (enc *keyValueEncoder) AppendObject(obj zapcore.ObjectMarshaler) error {
	enc.buf.AppendByte('{')
	err := obj.MarshalLogObject(enc)
	enc.buf.AppendByte('}')
	return err
}

func (enc *keyValueEncoder) AppendBool(val bool) {
	enc.addElementSeparator()
	enc.buf.AppendBool(val)
}

func (enc *keyValueEncoder) AppendByteString(val []byte) {
	enc.buf.Write(val)
}

// appendComplex appends the encoded form of the provided complex128 value.
// precision specifies the encoding precision for the real and imaginary
// components of the complex number.
func (enc *keyValueEncoder) appendComplex(val complex128, precision int) {
	enc.addElementSeparator()
	// Cast to a platform-independent, fixed-size type.
	r, i := float64(real(val)), float64(imag(val))
	enc.buf.AppendByte('"')
	// Because we're always in a quoted string, we can use strconv without
	// special-casing NaN and +/-Inf.
	enc.buf.AppendFloat(r, precision)
	// If imaginary part is less than 0, minus (-) sign is added by default
	// by AppendFloat.
	if i >= 0 {
		enc.buf.AppendByte('+')
	}
	enc.buf.AppendFloat(i, precision)
	enc.buf.AppendByte('i')
	enc.buf.AppendByte('"')
}

func (enc *keyValueEncoder) AppendDuration(val time.Duration) {
	cur := enc.buf.Len()
	if e := enc.EncodeDuration; e != nil {
		e(val, enc)
	}
	if cur == enc.buf.Len() {
		// User-supplied EncodeDuration is a no-op. Fall back to nanoseconds to keep
		// JSON valid.
		enc.AppendInt64(int64(val))
	}
}

func (enc *keyValueEncoder) AppendInt64(val int64) {
	//enc.addElementSeparator()
	enc.buf.AppendInt(val)
}

func (enc *keyValueEncoder) AppendReflected(val interface{}) error {
	valueBytes, err := enc.encodeReflected(val)
	if err != nil {
		return err
	}
	enc.addElementSeparator()
	_, err = enc.buf.Write(valueBytes)
	return err
}

func (enc *keyValueEncoder) AppendString(val string) {
	enc.buf.AppendString(val)
}

func (enc *keyValueEncoder) AppendTimeLayout(time time.Time, layout string) {
	enc.buf.AppendTime(time, layout)
}

func (enc *keyValueEncoder) AppendTime(val time.Time) {
	cur := enc.buf.Len()
	if e := enc.EncodeTime; e != nil {
		e(val, enc)
	}
	if cur == enc.buf.Len() {
		enc.AppendInt64(val.UnixNano())
	}
}

func (enc *keyValueEncoder) AppendUint64(val uint64) {
	enc.addElementSeparator()
	enc.buf.AppendUint(val)
}

func (enc *keyValueEncoder) appendFloat(val float64, bitSize int) {
	enc.addElementSeparator()
	switch {
	case math.IsNaN(val):
		enc.buf.AppendString(`"NaN"`)
	case math.IsInf(val, 1):
		enc.buf.AppendString(`"+Inf"`)
	case math.IsInf(val, -1):
		enc.buf.AppendString(`"-Inf"`)
	default:
		enc.buf.AppendFloat(val, bitSize)
	}
}

func (enc *keyValueEncoder) AddInt(k string, v int)         { enc.AddInt64(k, int64(v)) }
func (enc *keyValueEncoder) AddInt32(k string, v int32)     { enc.AddInt64(k, int64(v)) }
func (enc *keyValueEncoder) AddInt16(k string, v int16)     { enc.AddInt64(k, int64(v)) }
func (enc *keyValueEncoder) AddInt8(k string, v int8)       { enc.AddInt64(k, int64(v)) }
func (enc *keyValueEncoder) AddUint(k string, v uint)       { enc.AddUint64(k, uint64(v)) }
func (enc *keyValueEncoder) AddUint32(k string, v uint32)   { enc.AddUint64(k, uint64(v)) }
func (enc *keyValueEncoder) AddUint16(k string, v uint16)   { enc.AddUint64(k, uint64(v)) }
func (enc *keyValueEncoder) AddUint8(k string, v uint8)     { enc.AddUint64(k, uint64(v)) }
func (enc *keyValueEncoder) AddUintptr(k string, v uintptr) { enc.AddUint64(k, uint64(v)) }
func (enc *keyValueEncoder) AppendComplex64(v complex64)    { enc.appendComplex(complex128(v), 32) }
func (enc *keyValueEncoder) AppendComplex128(v complex128)  { enc.appendComplex(complex128(v), 64) }
func (enc *keyValueEncoder) AppendFloat64(v float64)        { enc.appendFloat(v, 64) }
func (enc *keyValueEncoder) AppendFloat32(v float32)        { enc.appendFloat(float64(v), 32) }
func (enc *keyValueEncoder) AppendInt(v int)                { enc.AppendInt64(int64(v)) }
func (enc *keyValueEncoder) AppendInt32(v int32)            { enc.AppendInt64(int64(v)) }
func (enc *keyValueEncoder) AppendInt16(v int16)            { enc.AppendInt64(int64(v)) }
func (enc *keyValueEncoder) AppendInt8(v int8)              { enc.AppendInt64(int64(v)) }
func (enc *keyValueEncoder) AppendUint(v uint)              { enc.AppendUint64(uint64(v)) }
func (enc *keyValueEncoder) AppendUint32(v uint32)          { enc.AppendUint64(uint64(v)) }
func (enc *keyValueEncoder) AppendUint16(v uint16)          { enc.AppendUint64(uint64(v)) }
func (enc *keyValueEncoder) AppendUint8(v uint8)            { enc.AppendUint64(uint64(v)) }
func (enc *keyValueEncoder) AppendUintptr(v uintptr)        { enc.AppendUint64(uint64(v)) }

// key 和value之间的分隔符
func (enc *keyValueEncoder) addKey(key string) {
	enc.addElementSeparator()
	enc.buf.AppendString(key)
	enc.buf.AppendByte('=')
}

// 添加列与列之间的分隔符
func (enc *keyValueEncoder) addElementSeparator() {
	enc.buf.AppendByte('\t')
}

func (enc *keyValueEncoder) Clone() zapcore.Encoder {
	clone := enc.clone()
	// 将输入写入clone的切片中
	clone.buf.Write(enc.buf.Bytes())
	return clone
}

// 从对象池中获取对象
func getKeyValueEncoder() *keyValueEncoder {
	return _keyValueEncoderPool.Get().(*keyValueEncoder)
}

// 将对象还给对象池
func putJSONEncoder(enc *keyValueEncoder) {
	if enc.reflectBuf != nil {
		enc.reflectBuf.Free()
	}
	enc.EncoderConfig = nil
	enc.buf = nil
	enc.reflectBuf = nil
	enc.reflectEnc = nil
	_keyValueEncoderPool.Put(enc)
}

func (enc *keyValueEncoder) clone() *keyValueEncoder {
	// 从对象池中获取一个编码器
	clone := getKeyValueEncoder()
	clone.EncoderConfig = enc.EncoderConfig
	// 获取字节切片
	clone.buf = _bufferPool.Get()
	return clone
}

func (enc *keyValueEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	clone := enc.clone()
	//buf := clone.buf

	if clone.TimeKey != "" {
		clone.buf.AppendString(enc.TimeKey)
		clone.buf.AppendByte('=')
		clone.EncodeTime(ent.Time, clone)
		clone.buf.AppendByte('\t')
	}

	// 时间
	//if !ent.Time.IsZero() {
	//	clone.buf.AppendString(fmt.Sprintf("ts=%s ", ent.Time.Format("2006-01-02T15:04:05.000Z0700")))
	//}

	// foo.go:123
	if clone.CallerKey != "" && ent.Caller.Defined {
		clone.buf.AppendString(enc.CallerKey)
		clone.buf.AppendByte('=')
		clone.EncodeCaller(ent.Caller, clone)
		clone.buf.AppendByte('\t')
	}

	if clone.StacktraceKey != "" {
		clone.buf.AppendString(fmt.Sprintf("%v=%v", clone.StacktraceKey, ent.Stack))
	}
	clone.buf.AppendByte('\t')
	//// 日志级别
	//clone.buf.AppendString(fmt.Sprintf("logLev=[%s]", ent.Level.CapitalString()))

	if clone.MessageKey != "" {
		clone.buf.AppendString(fmt.Sprintf("%v=%v", clone.MessageKey, ent.Message))
	}

	clone.buf.AppendByte('\t')

	for i := range fields {
		fields[i].AddTo(clone)
	}

	// 遍历并格式化字段
	//for _, field := range fields {
	//	clone.buf.AppendString(fmt.Sprintf("%s=%v ", field.Key, field.Interface))
	//}
	//clone.buf.AppendString(enc.LineEnding)
	//// 换行
	clone.buf.AppendString("\n")

	ret := clone.buf
	putJSONEncoder(clone)
	return ret, nil
}

func newKeyValueEncoder(config *zapcore.EncoderConfig) zapcore.Encoder {
	return &keyValueEncoder{
		EncoderConfig: config,
		buf:           _bufferPool.Get(),
	}
}
