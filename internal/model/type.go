package model

import (
	"database/sql/driver"
	"fmt"
	"math"
	"strconv"
	"time"
)

// MyString defines SQL string type.
//
// Note: unable to write NULL value to database; read NULL value from database as "".
type MyString string

const MyStringNull MyString = ""

// Value implements the driver.Valuer interface, which is automatically called when writing to the database.
func (s MyString) Value() (driver.Value, error) {
	return s.String(), nil
}

// Scan implements sql.Scanner interface, which is automatically called when reading from the database.
func (s *MyString) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		*s = MyString(v)
	case []byte:
		*s = MyString(v)
	case nil:
		*s = MyStringNull
	default:
		return fmt.Errorf("unexpected type for MyString: %T", src)
	}
	return nil
}

// DateString returns SQL date format string.
func (s MyString) DateString() (resDate MyDateString) {
	_, err := time.Parse("2006-01-02", s.String())
	if err == nil {
		resDate = MyDateString(s)
	}
	return
}

// String returns string type value.
func (s MyString) String() string {
	return string(s)
}

// MyDateString defines SQL date format string type.
//
// Note: "" as NULL value.
type MyDateString string

const MyDateStringNull MyDateString = ""

// Value implements the driver.Valuer interface, which is automatically called when writing to the database.
func (d MyDateString) Value() (driver.Value, error) {
	if d == MyDateStringNull {
		//* return NULL
		return nil, nil
	} else {
		_, err := time.Parse("2006-01-02", d.String())
		if err == nil {
			return d.String(), nil
		} else {
			//? return NULL if it does not match the format
			return nil, nil
		}
	}
}

// Scan implements sql.Scanner interface, which is automatically called when reading from the database.
func (d *MyDateString) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		*d = MyDateString(v)
	case []byte:
		*d = MyDateString(v)
	case time.Time:
		*d = MyDateString(v.Format("2006-01-02"))
	case nil:
		*d = MyDateStringNull
	default:
		return fmt.Errorf("unexpected type for MyDateString: %T", src)
	}
	return nil
}

// Valid returns true when the date string is valid.
func (d MyDateString) Valid() bool {
	res, _ := d.Value()
	return res != nil
}

// String returns string type value.
func (d MyDateString) String() string {
	return string(d)
}

// MyTimeString defines SQL datetime format string type.
//
// Note: "" as NULL value.
type MyDatetimeString string

const MyDatetimeStringNull MyDatetimeString = ""

// Value implements the driver.Valuer interface, which is automatically called when writing to the database.
func (d MyDatetimeString) Value() (driver.Value, error) {
	if d == MyDatetimeStringNull {
		//* return NULL
		return nil, nil
	} else {
		_, err := time.Parse("2006-01-02 15:04:05", d.String())
		if err == nil {
			return d.String(), nil
		} else {
			//? return NULL if it does not match the format
			return nil, nil
		}
	}
}

// Scan implements sql.Scanner interface, which is automatically called when reading from the database.
func (d *MyDatetimeString) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		*d = MyDatetimeString(v)
	case []byte:
		*d = MyDatetimeString(v)
	case time.Time:
		*d = MyDatetimeString(v.Format("2006-01-02 15:04:05"))
	case nil:
		*d = MyDatetimeStringNull
	default:
		return fmt.Errorf("unexpected type for MyDatetimeString: %T", src)
	}
	return nil
}

// Valid returns true when the datetime string is valid.
func (d MyDatetimeString) Valid() bool {
	res, _ := d.Value()
	return res != nil
}

// String returns string type value.
func (d MyDatetimeString) String() string {
	return string(d)
}

// MyInt defines SQL int type
//
// Note: unable to write NULL value to database; read NULL value from database as 0.
type MyInt int

const MyIntNull MyInt = 0

// Value implements the driver.Valuer interface, which is automatically called when writing to the database.
func (i MyInt) Value() (driver.Value, error) {
	//! driver.value need `int64` type param
	return i.Int64(), nil
}

// Scan implements sql.Scanner interface, which is automatically called when reading from the database.
func (i *MyInt) Scan(src interface{}) error {
	switch v := src.(type) {
	case int:
		*i = MyInt(v)
	case int64:
		*i = MyInt(v)
	case float64:
		*i = MyInt(math.Round(v))
	case []byte: //? maybe can scan the []byte([]uint8) type
		strValue := string(v)
		intValue, err := strconv.ParseInt(strValue, 10, 32)
		if err != nil {
			return fmt.Errorf("unexpected type for MyInt: %T", src)
		}
		*i = MyInt(intValue)
	case nil:
		*i = MyIntNull
	default:
		return fmt.Errorf("unexpected type for MyInt: %T", src)
	}
	return nil
}

// Int returns int type value.
func (i MyInt) Int() int {
	return int(i)
}

// Int64 returns int64 type value.
func (i MyInt) Int64() int64 {
	return int64(i)
}

// String returns string type value.
func (i MyInt) String() string {
	return strconv.Itoa(i.Int())
}

// MyInt64 defines SQL int64 type
//
// Note: unable to write NULL value to database; read NULL value from database as 0.
type MyInt64 int64

const MyInt64Null MyInt64 = 0

// Value implements the driver.Valuer interface, which is automatically called when writing to the database.
func (i MyInt64) Value() (driver.Value, error) {
	return i.Int64(), nil
}

// Scan implements sql.Scanner interface, which is automatically called when reading from the database.
func (i *MyInt64) Scan(src interface{}) error {
	switch v := src.(type) {
	case int64:
		*i = MyInt64(v)
	case float64:
		*i = MyInt64(math.Round(v))
	case []byte: //? maybe can scan the []byte([]uint8) type
		strValue := string(v)
		intValue, err := strconv.ParseInt(strValue, 10, 64)
		if err != nil {
			return fmt.Errorf("unexpected type for MyInt64: %T", src)
		}
		*i = MyInt64(intValue)
	case nil:
		*i = MyInt64Null
	default:
		return fmt.Errorf("unexpected type for MyInt64: %T", src)
	}
	return nil
}

// Int64 returns int64 type value.
func (i MyInt64) Int64() int64 {
	return int64(i)
}

// String returns string type value.
func (i MyInt64) String() string {
	return strconv.FormatInt(i.Int64(), 10)
}

// MyInt16 defines SQL int16 type
//
// Note: unable to write NULL value to database; read NULL value from database as 0.
type MyInt16 int16

const MyInt16Null MyInt16 = 0

// Value implements the driver.Valuer interface, which is automatically called when writing to the database.
func (i MyInt16) Value() (driver.Value, error) {
	//! driver.value need `int64` type param
	return i.Int64(), nil
}

// Scan implements sql.Scanner interface, which is automatically called when reading from the database.
func (i *MyInt16) Scan(src interface{}) error {
	switch v := src.(type) {
	case int64: //* accept int64 type
		*i = MyInt16(v)
	case int16: // the theory will not return to this type, only reserved
		*i = MyInt16(v)
	case float64:
		*i = MyInt16(math.Round(v))
	case []byte: //? maybe can scan the []byte([]uint8) type
		strValue := string(v)
		intValue, err := strconv.ParseInt(strValue, 10, 16)
		if err != nil {
			return fmt.Errorf("unexpected type for MyInt16: %T", src)
		}
		*i = MyInt16(intValue)
	case nil:
		*i = MyInt16Null
	default:
		return fmt.Errorf("unexpected type for MyInt16: %T", src)
	}
	return nil
}

// Int16 returns int16 type value.
func (i MyInt16) Int16() int16 {
	return int16(i)
}

// Int64 returns int64 type value.
func (i MyInt16) Int64() int64 {
	return int64(i)
}

// String returns string type value.
func (i MyInt16) String() string {
	return strconv.FormatInt(i.Int64(), 10)
}

// MyFloat64 defines SQL float64 type
//
// Note: unable to write NULL value to database; read NULL value from database as 0.0.
type MyFloat64 float64

const MyFloat64Null MyFloat64 = 0.0

// Value implements the driver.Valuer interface, which is automatically called when writing to the database.
func (f MyFloat64) Value() (driver.Value, error) {
	return f.Float64(), nil
}

// Scan implements sql.Scanner interface, which is automatically called when reading from the database.
func (f *MyFloat64) Scan(src interface{}) error {
	switch v := src.(type) {
	case float64:
		*f = MyFloat64(v)
	case float32:
		*f = MyFloat64(v)
	case []byte: //? maybe can scan the []byte([]uint8) type
		strValue := string(v)
		floatValue, err := strconv.ParseFloat(strValue, 64)
		if err != nil {
			return fmt.Errorf("unexpected type for MyFloat64: %T", src)
		}
		*f = MyFloat64(floatValue)
	case nil:
		*f = MyFloat64Null
	default:
		return fmt.Errorf("unexpected type for MyFloat64: %T", src)
	}
	return nil
}

// Float64 returns float64 type value.
func (f MyFloat64) Float64() float64 {
	return float64(f)
}

// String returns string type value.
func (f MyFloat64) String() string {
	return strconv.FormatFloat(f.Float64(), 'f', -1, 64)
}

// MyFloat32 defines SQL float32 type
//
// Note: unable to write NULL value to database; read NULL value from database as 0.0.
type MyFloat32 float32

const MyFloat32Null MyFloat32 = 0.0

// Value implements the driver.Valuer interface, which is automatically called when writing to the database.
func (f MyFloat32) Value() (driver.Value, error) {
	return f.Float64(), nil
}

// Scan implements sql.Scanner interface, which is automatically called when reading from the database.
func (f *MyFloat32) Scan(src interface{}) error {
	switch v := src.(type) {
	case float32:
		*f = MyFloat32(v)
	case float64:
		*f = MyFloat32(v)
	case []byte: //? maybe can scan the []byte([]uint8) type
		strValue := string(v)
		floatValue, err := strconv.ParseFloat(strValue, 32)
		if err != nil {
			return fmt.Errorf("unexpected type for MyFloat32: %T", src)
		}
		*f = MyFloat32(floatValue)
	case nil:
		*f = MyFloat32Null
	default:
		return fmt.Errorf("unexpected type for MyFloat32: %T", src)
	}
	return nil
}

// Float32 returns float32 type value.
func (f MyFloat32) Float32() float32 {
	return float32(f)
}

// Float64 returns float64 type value.
func (f MyFloat32) Float64() float64 {
	return float64(f)
}

// String returns string type value.
func (f MyFloat32) String() string {
	return strconv.FormatFloat(f.Float64(), 'f', -1, 32)
}

// MyBool defines SQL bool type
//
// Note: unable to write NULL value to database; read NULL value from database as false.
type MyBool bool

const MyBoolNull MyBool = false

// Value implements the driver.Valuer interface, which is automatically called when writing to the database.
func (b MyBool) Value() (driver.Value, error) {
	return b.Bool(), nil
}

// Scan implements sql.Scanner interface, which is automatically called when reading from the database.
func (b *MyBool) Scan(src interface{}) error {
	switch v := src.(type) {
	case bool:
		*b = MyBool(v)
	case int64: //? maybe can scan the int64 type
		*b = MyBool(v == 1)
	case nil:
		*b = MyBoolNull
	default:
		return fmt.Errorf("unexpected type for MyBool: %T", src)
	}
	return nil
}

// Bool returns bool type value.
func (b MyBool) Bool() bool {
	return bool(b)
}

// NullInt64 defines SQL int64 or NULL type
//
// Note: 0 as NULL
type NullInt64 int64

const NullInt64Null NullInt64 = 0

// Value implements the driver.Valuer interface, which is automatically called when writing to the database.
func (i NullInt64) Value() (driver.Value, error) {
	if i == NullInt64Null {
		return nil, nil
	}
	return i.Int64(), nil
}

// Scan implements sql.Scanner interface, which is automatically called when reading from the database.
func (i *NullInt64) Scan(src interface{}) error {
	switch v := src.(type) {
	case int64:
		*i = NullInt64(v)
	case float64:
		*i = NullInt64(math.Round(v))
	case []byte: //? maybe can scan the []byte([]uint8) type
		strValue := string(v)
		intValue, err := strconv.ParseInt(strValue, 10, 64)
		if err != nil {
			return fmt.Errorf("unexpected type for NullInt64: %T", src)
		}
		*i = NullInt64(intValue)
	case nil:
		*i = NullInt64Null
	default:
		return fmt.Errorf("unexpected type for NullInt64: %T", src)
	}
	return nil
}

// Int64 returns int64 type value.
func (i NullInt64) Int64() int64 {
	return int64(i)
}

// String returns string type value.
func (i NullInt64) String() string {
	return strconv.FormatInt(i.Int64(), 10)
}
