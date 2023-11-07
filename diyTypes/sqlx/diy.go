package sqlx

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"github.com/shopspring/decimal"
	"strconv"
)

type NullString sql.NullString
type NullInt64 sql.NullInt64
type NullInt16 sql.NullInt16
type NullInt32 sql.NullInt32
type NullTime sql.NullTime
type NullBool sql.NullBool
type NullByte sql.NullByte
type NullFloat64 sql.NullFloat64
type Decimal decimal.Decimal
type NullDecimal decimal.NullDecimal

func (rec *NullString) Scan(value interface{}) error {
	var x sql.NullString
	err := x.Scan(value)
	if err != nil {
		return err
	}
	*rec = NullString(x)
	return nil
}
func (rec NullString) Value() (driver.Value, error) {
	return sql.NullString(rec).Value()
}
func (rec *NullInt64) Scan(value interface{}) error {
	var x sql.NullInt64
	err := x.Scan(value)
	if err != nil {
		return err
	}
	*rec = NullInt64(x)
	return nil
}
func (rec NullInt64) Value() (driver.Value, error) {
	return sql.NullInt64(rec).Value()
}
func (rec *NullInt16) Scan(value interface{}) error {
	var x sql.NullInt16
	err := x.Scan(value)
	if err != nil {
		return err
	}
	*rec = NullInt16(x)
	return nil
}
func (rec NullInt16) Value() (driver.Value, error) {
	return sql.NullInt16(rec).Value()
}

func (rec *NullInt32) Scan(value interface{}) error {
	var x sql.NullInt32
	err := x.Scan(value)
	if err != nil {
		return err
	}
	*rec = NullInt32(x)
	return nil
}
func (rec NullInt32) Value() (driver.Value, error) {
	return sql.NullInt32(rec).Value()
}

func (rec *NullTime) Scan(value interface{}) error {
	var x sql.NullTime
	err := x.Scan(value)
	if err != nil {
		return err
	}
	*rec = NullTime(x)
	return nil
}
func (rec NullTime) Value() (driver.Value, error) {
	return sql.NullTime(rec).Value()
}

func (rec *NullBool) Scan(value interface{}) error {
	var x sql.NullBool
	err := x.Scan(value)
	if err != nil {
		return err
	}
	*rec = NullBool(x)
	return nil
}
func (rec NullBool) Value() (driver.Value, error) {
	return sql.NullBool(rec).Value()
}

func (rec *NullByte) Scan(value interface{}) error {
	var x sql.NullByte
	err := x.Scan(value)
	if err != nil {
		return err
	}
	*rec = NullByte(x)
	return nil
}
func (rec NullByte) Value() (driver.Value, error) {
	return sql.NullByte(rec).Value()
}

func (rec *NullFloat64) Scan(value interface{}) error {
	var x sql.NullFloat64
	err := x.Scan(value)
	if err != nil {
		return err
	}
	*rec = NullFloat64(x)
	return nil
}
func (rec NullFloat64) Value() (driver.Value, error) {
	return sql.NullFloat64(rec).Value()
}

func (d *NullDecimal) Scan(value interface{}) error {
	var x decimal.NullDecimal
	err := x.Scan(value)
	if err != nil {
		return err
	}
	*d = NullDecimal(x)
	return nil
}
func (d NullDecimal) Value() (driver.Value, error) {
	return decimal.NullDecimal(d).Value()
}

func (d *Decimal) Scan(value interface{}) error {
	var x decimal.Decimal
	err := x.Scan(value)
	if err != nil {
		return err
	}
	*d = Decimal(x)
	return nil
}
func (d Decimal) Value() (driver.Value, error) {
	return decimal.Decimal(d).String(), nil
}

func (d Decimal) MarshalJSON() ([]byte, error) {
	return json.Marshal(decimal.Decimal(d).InexactFloat64())
}
func (d *Decimal) UnmarshalJSON(data []byte) error {
	var t decimal.Decimal
	err := t.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	*d = Decimal(t)
	return nil
}

func (d NullDecimal) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return jsonNull, nil
	}
	return json.Marshal(d.Decimal.InexactFloat64())
}
func (d *NullDecimal) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		d.Valid = false
		return nil
	}
	err := d.Decimal.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	d.Valid = true
	return nil
}

func (receiver NullString) MarshalJSON() ([]byte, error) {
	if receiver.Valid {
		return json.Marshal(receiver.String)
	}
	return jsonNull, nil
}
func (receiver *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		receiver.Valid = false
		return nil
	}
	receiver.Valid = true
	receiver.String = string(data)
	return nil
}
func (receiver NullInt64) MarshalJSON() ([]byte, error) {
	if receiver.Valid {
		return json.Marshal(receiver.Int64)
	}
	return jsonNull, nil
}

func (receiver *NullInt64) UnmarshalJSON(data []byte) error {
	if !bytes.Equal(data, jsonNull) {
		var err error
		receiver.Int64, err = strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return err
		}
		receiver.Valid = true
	}
	return nil
}

func (receiver NullInt16) MarshalJSON() ([]byte, error) {
	if receiver.Valid {
		return json.Marshal(receiver.Int16)
	}
	return jsonNull, nil
}

func (receiver *NullInt16) UnmarshalJSON(data []byte) error {
	if !bytes.Equal(data, jsonNull) {
		in, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return err
		}
		receiver.Valid = true
		receiver.Int16 = int16(in)
	}
	return nil
}

func (receiver NullInt32) MarshalJSON() ([]byte, error) {
	if receiver.Valid {
		return json.Marshal(receiver.Int32)
	}
	return jsonNull, nil
}

func (receiver *NullInt32) UnmarshalJSON(data []byte) error {
	if !bytes.Equal(data, jsonNull) {
		in, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return err
		}
		receiver.Valid = true
		receiver.Int32 = int32(in)
	}
	return nil
}

func (receiver *NullByte) MarshalJSON() ([]byte, error) {
	if receiver.Valid {
		return json.Marshal(receiver.Byte)
	}
	return jsonNull, nil
}

func (receiver *NullByte) UnmarshalJSON(data []byte) error {
	if !bytes.Equal(data, jsonNull) {
		in, err := strconv.ParseUint(string(data), 10, 64)
		if err != nil {
			return err
		}
		receiver.Valid = true
		receiver.Byte = uint8(in)
	}
	return nil
}

func (receiver NullFloat64) MarshalJSON() ([]byte, error) {
	if receiver.Valid {
		return json.Marshal(receiver.Float64)
	}
	return jsonNull, nil
}

func (receiver *NullFloat64) UnmarshalJSON(data []byte) error {
	if !bytes.Equal(data, jsonNull) {
		in, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return err
		}
		receiver.Valid = true
		receiver.Float64 = in
	}
	return nil
}

func (receiver NullBool) MarshalJSON() ([]byte, error) {
	if receiver.Valid {
		return json.Marshal(receiver.Bool)
	}
	return jsonNull, nil
}

func (receiver *NullBool) UnmarshalJSON(data []byte) error {
	if !bytes.Equal(data, jsonNull) {
		in, err := strconv.ParseBool(string(data))
		if err != nil {
			return err
		}
		receiver.Valid = true
		receiver.Bool = in
	}
	return nil
}

func (receiver NullTime) MarshalJSON() ([]byte, error) {
	if receiver.Valid {
		return receiver.Time.MarshalJSON()
	}
	return jsonNull, nil
}

var jsonNull = []byte("null")

func (receiver *NullTime) UnmarshalJSON(data []byte) error {
	if !bytes.Equal(data, jsonNull) {
		receiver.Valid = true
		return receiver.Time.UnmarshalJSON(data)
	}
	return nil
}
