// DO NOT EDIT!
// Code generated by ffjson <https://github.com/pquerna/ffjson>
// source: topic_route.go
// DO NOT EDIT!

package rocketmq

import (
	"bytes"
	"fmt"

	fflib "github.com/pquerna/ffjson/fflib/v1"
)

func (mj *BrokerData) MarshalJSON() ([]byte, error) {
	var buf fflib.Buffer
	if mj == nil {
		buf.WriteString("null")
		return buf.Bytes(), nil
	}
	err := mj.MarshalJSONBuf(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (mj *BrokerData) MarshalJSONBuf(buf fflib.EncodingBuffer) error {
	if mj == nil {
		buf.WriteString("null")
		return nil
	}
	var err error
	var obj []byte
	_ = obj
	_ = err
	buf.WriteString(`{"brokerName":`)
	fflib.WriteJsonString(buf, string(mj.BrokerName))
	/* Falling back. type=map[int64]string kind=map */
	buf.WriteString(`,"brokerAddrs":`)
	err = buf.Encode(mj.BrokerAddrs)
	if err != nil {
		return err
	}
	buf.WriteByte('}')
	return nil
}

const (
	ffj_t_BrokerDatabase = iota
	ffj_t_BrokerDatano_such_key

	ffj_t_BrokerData_BrokerName

	ffj_t_BrokerData_BrokerAddrs
)

var ffj_key_BrokerData_BrokerName = []byte("brokerName")

var ffj_key_BrokerData_BrokerAddrs = []byte("brokerAddrs")

func (uj *BrokerData) UnmarshalJSON(input []byte) error {
	fs := fflib.NewFFLexer(input)
	return uj.UnmarshalJSONFFLexer(fs, fflib.FFParse_map_start)
}

func (uj *BrokerData) UnmarshalJSONFFLexer(fs *fflib.FFLexer, state fflib.FFParseState) error {
	var err error = nil
	currentKey := ffj_t_BrokerDatabase
	_ = currentKey
	tok := fflib.FFTok_init
	wantedTok := fflib.FFTok_init

mainparse:
	for {
		tok = fs.Scan()
		//	println(fmt.Sprintf("debug: tok: %v  state: %v", tok, state))
		if tok == fflib.FFTok_error {
			goto tokerror
		}

		switch state {

		case fflib.FFParse_map_start:
			if tok != fflib.FFTok_left_bracket {
				wantedTok = fflib.FFTok_left_bracket
				goto wrongtokenerror
			}
			state = fflib.FFParse_want_key
			continue

		case fflib.FFParse_after_value:
			if tok == fflib.FFTok_comma {
				state = fflib.FFParse_want_key
			} else if tok == fflib.FFTok_right_bracket {
				goto done
			} else {
				wantedTok = fflib.FFTok_comma
				goto wrongtokenerror
			}

		case fflib.FFParse_want_key:
			// json {} ended. goto exit. woo.
			if tok == fflib.FFTok_right_bracket {
				goto done
			}
			if tok != fflib.FFTok_string {
				wantedTok = fflib.FFTok_string
				goto wrongtokenerror
			}

			kn := fs.Output.Bytes()
			if len(kn) <= 0 {
				// "" case. hrm.
				currentKey = ffj_t_BrokerDatano_such_key
				state = fflib.FFParse_want_colon
				goto mainparse
			} else {
				switch kn[0] {

				case 'b':

					if bytes.Equal(ffj_key_BrokerData_BrokerName, kn) {
						currentKey = ffj_t_BrokerData_BrokerName
						state = fflib.FFParse_want_colon
						goto mainparse

					} else if bytes.Equal(ffj_key_BrokerData_BrokerAddrs, kn) {
						currentKey = ffj_t_BrokerData_BrokerAddrs
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				}

				if fflib.EqualFoldRight(ffj_key_BrokerData_BrokerAddrs, kn) {
					currentKey = ffj_t_BrokerData_BrokerAddrs
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffj_key_BrokerData_BrokerName, kn) {
					currentKey = ffj_t_BrokerData_BrokerName
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				currentKey = ffj_t_BrokerDatano_such_key
				state = fflib.FFParse_want_colon
				goto mainparse
			}

		case fflib.FFParse_want_colon:
			if tok != fflib.FFTok_colon {
				wantedTok = fflib.FFTok_colon
				goto wrongtokenerror
			}
			state = fflib.FFParse_want_value
			continue
		case fflib.FFParse_want_value:

			if tok == fflib.FFTok_left_brace || tok == fflib.FFTok_left_bracket || tok == fflib.FFTok_integer || tok == fflib.FFTok_double || tok == fflib.FFTok_string || tok == fflib.FFTok_bool || tok == fflib.FFTok_null {
				switch currentKey {

				case ffj_t_BrokerData_BrokerName:
					goto handle_BrokerName

				case ffj_t_BrokerData_BrokerAddrs:
					goto handle_BrokerAddrs

				case ffj_t_BrokerDatano_such_key:
					err = fs.SkipField(tok)
					if err != nil {
						return fs.WrapErr(err)
					}
					state = fflib.FFParse_after_value
					goto mainparse
				}
			} else {
				goto wantedvalue
			}
		}
	}

handle_BrokerName:

	/* handler: uj.BrokerName type=string kind=string quoted=false*/

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			outBuf := fs.Output.Bytes()

			uj.BrokerName = string(string(outBuf))

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_BrokerAddrs:

	/* handler: uj.BrokerAddrs type=map[int64]string kind=map quoted=false*/

	{

		{
			if tok != fflib.FFTok_left_bracket && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for ", tok))
			}
		}

		if tok == fflib.FFTok_null {
			uj.BrokerAddrs = nil
		} else {

			uj.BrokerAddrs = make(map[int64]string, 0)

			wantVal := true

			for {

				var k int64

				var tmp_uj__BrokerAddrs string

				tok = fs.Scan()
				if tok == fflib.FFTok_error {
					goto tokerror
				}
				if tok == fflib.FFTok_right_bracket {
					break
				}

				if tok == fflib.FFTok_comma {
					if wantVal == true {
						// TODO(pquerna): this isn't an ideal error message, this handles
						// things like [,,,] as an array value.
						return fs.WrapErr(fmt.Errorf("wanted value token, but got token: %v", tok))
					}
					continue
				} else {
					wantVal = true
				}

				/* handler: k type=int64 kind=int64 quoted=false*/

				{
					if tok != fflib.FFTok_integer && tok != fflib.FFTok_null {
						return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for int64", tok))
					}
				}

				{

					if tok == fflib.FFTok_null {

					} else {

						tval, err := fflib.ParseInt(fs.Output.Bytes(), 10, 64)

						if err != nil {
							return fs.WrapErr(err)
						}

						k = int64(tval)

					}
				}

				// Expect ':' after key
				tok = fs.Scan()
				if tok != fflib.FFTok_colon {
					return fs.WrapErr(fmt.Errorf("wanted colon token, but got token: %v", tok))
				}

				tok = fs.Scan()
				/* handler: tmp_uj__BrokerAddrs type=string kind=string quoted=false*/

				{

					{
						if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
							return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
						}
					}

					if tok == fflib.FFTok_null {

					} else {

						outBuf := fs.Output.Bytes()

						tmp_uj__BrokerAddrs = string(string(outBuf))

					}
				}

				uj.BrokerAddrs[k] = tmp_uj__BrokerAddrs

				wantVal = false
			}

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

wantedvalue:
	return fs.WrapErr(fmt.Errorf("wanted value token, but got token: %v", tok))
wrongtokenerror:
	return fs.WrapErr(fmt.Errorf("ffjson: wanted token: %v, but got token: %v output=%s", wantedTok, tok, fs.Output.String()))
tokerror:
	if fs.BigError != nil {
		return fs.WrapErr(fs.BigError)
	}
	err = fs.Error.ToError()
	if err != nil {
		return fs.WrapErr(err)
	}
	panic("ffjson-generated: unreachable, please report bug.")
done:

	return nil
}

func (mj *QueueData) MarshalJSON() ([]byte, error) {
	var buf fflib.Buffer
	if mj == nil {
		buf.WriteString("null")
		return buf.Bytes(), nil
	}
	err := mj.MarshalJSONBuf(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (mj *QueueData) MarshalJSONBuf(buf fflib.EncodingBuffer) error {
	if mj == nil {
		buf.WriteString("null")
		return nil
	}
	var err error
	var obj []byte
	_ = obj
	_ = err
	buf.WriteString(`{"BrokerName":`)
	fflib.WriteJsonString(buf, string(mj.BrokerName))
	buf.WriteString(`,"ReadQueueNums":`)
	fflib.FormatBits2(buf, uint64(mj.ReadQueueNums), 10, mj.ReadQueueNums < 0)
	buf.WriteString(`,"WriteQueueNums":`)
	fflib.FormatBits2(buf, uint64(mj.WriteQueueNums), 10, mj.WriteQueueNums < 0)
	buf.WriteString(`,"Perm":`)
	fflib.FormatBits2(buf, uint64(mj.Perm), 10, mj.Perm < 0)
	buf.WriteString(`,"TopicSynFlag":`)
	fflib.FormatBits2(buf, uint64(mj.TopicSynFlag), 10, mj.TopicSynFlag < 0)
	buf.WriteByte('}')
	return nil
}

const (
	ffj_t_QueueDatabase = iota
	ffj_t_QueueDatano_such_key

	ffj_t_QueueData_BrokerName

	ffj_t_QueueData_ReadQueueNums

	ffj_t_QueueData_WriteQueueNums

	ffj_t_QueueData_Perm

	ffj_t_QueueData_TopicSynFlag
)

var ffj_key_QueueData_BrokerName = []byte("BrokerName")

var ffj_key_QueueData_ReadQueueNums = []byte("ReadQueueNums")

var ffj_key_QueueData_WriteQueueNums = []byte("WriteQueueNums")

var ffj_key_QueueData_Perm = []byte("Perm")

var ffj_key_QueueData_TopicSynFlag = []byte("TopicSynFlag")

func (uj *QueueData) UnmarshalJSON(input []byte) error {
	fs := fflib.NewFFLexer(input)
	return uj.UnmarshalJSONFFLexer(fs, fflib.FFParse_map_start)
}

func (uj *QueueData) UnmarshalJSONFFLexer(fs *fflib.FFLexer, state fflib.FFParseState) error {
	var err error = nil
	currentKey := ffj_t_QueueDatabase
	_ = currentKey
	tok := fflib.FFTok_init
	wantedTok := fflib.FFTok_init

mainparse:
	for {
		tok = fs.Scan()
		//	println(fmt.Sprintf("debug: tok: %v  state: %v", tok, state))
		if tok == fflib.FFTok_error {
			goto tokerror
		}

		switch state {

		case fflib.FFParse_map_start:
			if tok != fflib.FFTok_left_bracket {
				wantedTok = fflib.FFTok_left_bracket
				goto wrongtokenerror
			}
			state = fflib.FFParse_want_key
			continue

		case fflib.FFParse_after_value:
			if tok == fflib.FFTok_comma {
				state = fflib.FFParse_want_key
			} else if tok == fflib.FFTok_right_bracket {
				goto done
			} else {
				wantedTok = fflib.FFTok_comma
				goto wrongtokenerror
			}

		case fflib.FFParse_want_key:
			// json {} ended. goto exit. woo.
			if tok == fflib.FFTok_right_bracket {
				goto done
			}
			if tok != fflib.FFTok_string {
				wantedTok = fflib.FFTok_string
				goto wrongtokenerror
			}

			kn := fs.Output.Bytes()
			if len(kn) <= 0 {
				// "" case. hrm.
				currentKey = ffj_t_QueueDatano_such_key
				state = fflib.FFParse_want_colon
				goto mainparse
			} else {
				switch kn[0] {

				case 'B':

					if bytes.Equal(ffj_key_QueueData_BrokerName, kn) {
						currentKey = ffj_t_QueueData_BrokerName
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'P':

					if bytes.Equal(ffj_key_QueueData_Perm, kn) {
						currentKey = ffj_t_QueueData_Perm
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'R':

					if bytes.Equal(ffj_key_QueueData_ReadQueueNums, kn) {
						currentKey = ffj_t_QueueData_ReadQueueNums
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'T':

					if bytes.Equal(ffj_key_QueueData_TopicSynFlag, kn) {
						currentKey = ffj_t_QueueData_TopicSynFlag
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'W':

					if bytes.Equal(ffj_key_QueueData_WriteQueueNums, kn) {
						currentKey = ffj_t_QueueData_WriteQueueNums
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				}

				if fflib.EqualFoldRight(ffj_key_QueueData_TopicSynFlag, kn) {
					currentKey = ffj_t_QueueData_TopicSynFlag
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.SimpleLetterEqualFold(ffj_key_QueueData_Perm, kn) {
					currentKey = ffj_t_QueueData_Perm
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffj_key_QueueData_WriteQueueNums, kn) {
					currentKey = ffj_t_QueueData_WriteQueueNums
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffj_key_QueueData_ReadQueueNums, kn) {
					currentKey = ffj_t_QueueData_ReadQueueNums
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffj_key_QueueData_BrokerName, kn) {
					currentKey = ffj_t_QueueData_BrokerName
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				currentKey = ffj_t_QueueDatano_such_key
				state = fflib.FFParse_want_colon
				goto mainparse
			}

		case fflib.FFParse_want_colon:
			if tok != fflib.FFTok_colon {
				wantedTok = fflib.FFTok_colon
				goto wrongtokenerror
			}
			state = fflib.FFParse_want_value
			continue
		case fflib.FFParse_want_value:

			if tok == fflib.FFTok_left_brace || tok == fflib.FFTok_left_bracket || tok == fflib.FFTok_integer || tok == fflib.FFTok_double || tok == fflib.FFTok_string || tok == fflib.FFTok_bool || tok == fflib.FFTok_null {
				switch currentKey {

				case ffj_t_QueueData_BrokerName:
					goto handle_BrokerName

				case ffj_t_QueueData_ReadQueueNums:
					goto handle_ReadQueueNums

				case ffj_t_QueueData_WriteQueueNums:
					goto handle_WriteQueueNums

				case ffj_t_QueueData_Perm:
					goto handle_Perm

				case ffj_t_QueueData_TopicSynFlag:
					goto handle_TopicSynFlag

				case ffj_t_QueueDatano_such_key:
					err = fs.SkipField(tok)
					if err != nil {
						return fs.WrapErr(err)
					}
					state = fflib.FFParse_after_value
					goto mainparse
				}
			} else {
				goto wantedvalue
			}
		}
	}

handle_BrokerName:

	/* handler: uj.BrokerName type=string kind=string quoted=false*/

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			outBuf := fs.Output.Bytes()

			uj.BrokerName = string(string(outBuf))

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_ReadQueueNums:

	/* handler: uj.ReadQueueNums type=int32 kind=int32 quoted=false*/

	{
		if tok != fflib.FFTok_integer && tok != fflib.FFTok_null {
			return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for int32", tok))
		}
	}

	{

		if tok == fflib.FFTok_null {

		} else {

			tval, err := fflib.ParseInt(fs.Output.Bytes(), 10, 32)

			if err != nil {
				return fs.WrapErr(err)
			}

			uj.ReadQueueNums = int32(tval)

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_WriteQueueNums:

	/* handler: uj.WriteQueueNums type=int32 kind=int32 quoted=false*/

	{
		if tok != fflib.FFTok_integer && tok != fflib.FFTok_null {
			return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for int32", tok))
		}
	}

	{

		if tok == fflib.FFTok_null {

		} else {

			tval, err := fflib.ParseInt(fs.Output.Bytes(), 10, 32)

			if err != nil {
				return fs.WrapErr(err)
			}

			uj.WriteQueueNums = int32(tval)

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_Perm:

	/* handler: uj.Perm type=int32 kind=int32 quoted=false*/

	{
		if tok != fflib.FFTok_integer && tok != fflib.FFTok_null {
			return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for int32", tok))
		}
	}

	{

		if tok == fflib.FFTok_null {

		} else {

			tval, err := fflib.ParseInt(fs.Output.Bytes(), 10, 32)

			if err != nil {
				return fs.WrapErr(err)
			}

			uj.Perm = int32(tval)

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_TopicSynFlag:

	/* handler: uj.TopicSynFlag type=int32 kind=int32 quoted=false*/

	{
		if tok != fflib.FFTok_integer && tok != fflib.FFTok_null {
			return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for int32", tok))
		}
	}

	{

		if tok == fflib.FFTok_null {

		} else {

			tval, err := fflib.ParseInt(fs.Output.Bytes(), 10, 32)

			if err != nil {
				return fs.WrapErr(err)
			}

			uj.TopicSynFlag = int32(tval)

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

wantedvalue:
	return fs.WrapErr(fmt.Errorf("wanted value token, but got token: %v", tok))
wrongtokenerror:
	return fs.WrapErr(fmt.Errorf("ffjson: wanted token: %v, but got token: %v output=%s", wantedTok, tok, fs.Output.String()))
tokerror:
	if fs.BigError != nil {
		return fs.WrapErr(fs.BigError)
	}
	err = fs.Error.ToError()
	if err != nil {
		return fs.WrapErr(err)
	}
	panic("ffjson-generated: unreachable, please report bug.")
done:

	return nil
}

func (mj *TopicRouteData) MarshalJSON() ([]byte, error) {
	var buf fflib.Buffer
	if mj == nil {
		buf.WriteString("null")
		return buf.Bytes(), nil
	}
	err := mj.MarshalJSONBuf(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (mj *TopicRouteData) MarshalJSONBuf(buf fflib.EncodingBuffer) error {
	if mj == nil {
		buf.WriteString("null")
		return nil
	}
	var err error
	var obj []byte
	_ = obj
	_ = err
	buf.WriteString(`{"OrderTopicConf":`)
	fflib.WriteJsonString(buf, string(mj.OrderTopicConf))
	buf.WriteString(`,"QueueDatas":`)
	if mj.QueueDatas != nil {
		buf.WriteString(`[`)
		for i, v := range mj.QueueDatas {
			if i != 0 {
				buf.WriteString(`,`)
			}

			{

				if v == nil {
					buf.WriteString("null")
					return nil
				}

				err = v.MarshalJSONBuf(buf)
				if err != nil {
					return err
				}

			}
		}
		buf.WriteString(`]`)
	} else {
		buf.WriteString(`null`)
	}
	buf.WriteString(`,"BrokerDatas":`)
	if mj.BrokerDatas != nil {
		buf.WriteString(`[`)
		for i, v := range mj.BrokerDatas {
			if i != 0 {
				buf.WriteString(`,`)
			}

			{

				if v == nil {
					buf.WriteString("null")
					return nil
				}

				err = v.MarshalJSONBuf(buf)
				if err != nil {
					return err
				}

			}
		}
		buf.WriteString(`]`)
	} else {
		buf.WriteString(`null`)
	}
	buf.WriteByte('}')
	return nil
}

const (
	ffj_t_TopicRouteDatabase = iota
	ffj_t_TopicRouteDatano_such_key

	ffj_t_TopicRouteData_OrderTopicConf

	ffj_t_TopicRouteData_QueueDatas

	ffj_t_TopicRouteData_BrokerDatas
)

var ffj_key_TopicRouteData_OrderTopicConf = []byte("OrderTopicConf")

var ffj_key_TopicRouteData_QueueDatas = []byte("QueueDatas")

var ffj_key_TopicRouteData_BrokerDatas = []byte("BrokerDatas")

func (uj *TopicRouteData) UnmarshalJSON(input []byte) error {
	fs := fflib.NewFFLexer(input)
	return uj.UnmarshalJSONFFLexer(fs, fflib.FFParse_map_start)
}

func (uj *TopicRouteData) UnmarshalJSONFFLexer(fs *fflib.FFLexer, state fflib.FFParseState) error {
	var err error = nil
	currentKey := ffj_t_TopicRouteDatabase
	_ = currentKey
	tok := fflib.FFTok_init
	wantedTok := fflib.FFTok_init

mainparse:
	for {
		tok = fs.Scan()
		//	println(fmt.Sprintf("debug: tok: %v  state: %v", tok, state))
		if tok == fflib.FFTok_error {
			goto tokerror
		}

		switch state {

		case fflib.FFParse_map_start:
			if tok != fflib.FFTok_left_bracket {
				wantedTok = fflib.FFTok_left_bracket
				goto wrongtokenerror
			}
			state = fflib.FFParse_want_key
			continue

		case fflib.FFParse_after_value:
			if tok == fflib.FFTok_comma {
				state = fflib.FFParse_want_key
			} else if tok == fflib.FFTok_right_bracket {
				goto done
			} else {
				wantedTok = fflib.FFTok_comma
				goto wrongtokenerror
			}

		case fflib.FFParse_want_key:
			// json {} ended. goto exit. woo.
			if tok == fflib.FFTok_right_bracket {
				goto done
			}
			if tok != fflib.FFTok_string {
				wantedTok = fflib.FFTok_string
				goto wrongtokenerror
			}

			kn := fs.Output.Bytes()
			if len(kn) <= 0 {
				// "" case. hrm.
				currentKey = ffj_t_TopicRouteDatano_such_key
				state = fflib.FFParse_want_colon
				goto mainparse
			} else {
				switch kn[0] {

				case 'B':

					if bytes.Equal(ffj_key_TopicRouteData_BrokerDatas, kn) {
						currentKey = ffj_t_TopicRouteData_BrokerDatas
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'O':

					if bytes.Equal(ffj_key_TopicRouteData_OrderTopicConf, kn) {
						currentKey = ffj_t_TopicRouteData_OrderTopicConf
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'Q':

					if bytes.Equal(ffj_key_TopicRouteData_QueueDatas, kn) {
						currentKey = ffj_t_TopicRouteData_QueueDatas
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				}

				if fflib.EqualFoldRight(ffj_key_TopicRouteData_BrokerDatas, kn) {
					currentKey = ffj_t_TopicRouteData_BrokerDatas
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffj_key_TopicRouteData_QueueDatas, kn) {
					currentKey = ffj_t_TopicRouteData_QueueDatas
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.SimpleLetterEqualFold(ffj_key_TopicRouteData_OrderTopicConf, kn) {
					currentKey = ffj_t_TopicRouteData_OrderTopicConf
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				currentKey = ffj_t_TopicRouteDatano_such_key
				state = fflib.FFParse_want_colon
				goto mainparse
			}

		case fflib.FFParse_want_colon:
			if tok != fflib.FFTok_colon {
				wantedTok = fflib.FFTok_colon
				goto wrongtokenerror
			}
			state = fflib.FFParse_want_value
			continue
		case fflib.FFParse_want_value:

			if tok == fflib.FFTok_left_brace || tok == fflib.FFTok_left_bracket || tok == fflib.FFTok_integer || tok == fflib.FFTok_double || tok == fflib.FFTok_string || tok == fflib.FFTok_bool || tok == fflib.FFTok_null {
				switch currentKey {

				case ffj_t_TopicRouteData_OrderTopicConf:
					goto handle_OrderTopicConf

				case ffj_t_TopicRouteData_QueueDatas:
					goto handle_QueueDatas

				case ffj_t_TopicRouteData_BrokerDatas:
					goto handle_BrokerDatas

				case ffj_t_TopicRouteDatano_such_key:
					err = fs.SkipField(tok)
					if err != nil {
						return fs.WrapErr(err)
					}
					state = fflib.FFParse_after_value
					goto mainparse
				}
			} else {
				goto wantedvalue
			}
		}
	}

handle_OrderTopicConf:

	/* handler: uj.OrderTopicConf type=string kind=string quoted=false*/

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			outBuf := fs.Output.Bytes()

			uj.OrderTopicConf = string(string(outBuf))

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_QueueDatas:

	/* handler: uj.QueueDatas type=[]*rocketmq.QueueData kind=slice quoted=false*/

	{

		{
			if tok != fflib.FFTok_left_brace && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for ", tok))
			}
		}

		if tok == fflib.FFTok_null {
			uj.QueueDatas = nil
		} else {

			uj.QueueDatas = []*QueueData{}

			wantVal := true

			for {

				var tmp_uj__QueueDatas *QueueData

				tok = fs.Scan()
				if tok == fflib.FFTok_error {
					goto tokerror
				}
				if tok == fflib.FFTok_right_brace {
					break
				}

				if tok == fflib.FFTok_comma {
					if wantVal == true {
						// TODO(pquerna): this isn't an ideal error message, this handles
						// things like [,,,] as an array value.
						return fs.WrapErr(fmt.Errorf("wanted value token, but got token: %v", tok))
					}
					continue
				} else {
					wantVal = true
				}

				/* handler: tmp_uj__QueueDatas type=*rocketmq.QueueData kind=ptr quoted=false*/

				{
					if tok == fflib.FFTok_null {

						tmp_uj__QueueDatas = nil

						state = fflib.FFParse_after_value
						goto mainparse
					}

					if tmp_uj__QueueDatas == nil {
						tmp_uj__QueueDatas = new(QueueData)
					}

					err = tmp_uj__QueueDatas.UnmarshalJSONFFLexer(fs, fflib.FFParse_want_key)
					if err != nil {
						return err
					}
					state = fflib.FFParse_after_value
				}

				uj.QueueDatas = append(uj.QueueDatas, tmp_uj__QueueDatas)

				wantVal = false
			}
		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_BrokerDatas:

	/* handler: uj.BrokerDatas type=[]*rocketmq.BrokerData kind=slice quoted=false*/

	{

		{
			if tok != fflib.FFTok_left_brace && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for ", tok))
			}
		}

		if tok == fflib.FFTok_null {
			uj.BrokerDatas = nil
		} else {

			uj.BrokerDatas = []*BrokerData{}

			wantVal := true

			for {

				var tmp_uj__BrokerDatas *BrokerData

				tok = fs.Scan()
				if tok == fflib.FFTok_error {
					goto tokerror
				}
				if tok == fflib.FFTok_right_brace {
					break
				}

				if tok == fflib.FFTok_comma {
					if wantVal == true {
						// TODO(pquerna): this isn't an ideal error message, this handles
						// things like [,,,] as an array value.
						return fs.WrapErr(fmt.Errorf("wanted value token, but got token: %v", tok))
					}
					continue
				} else {
					wantVal = true
				}

				/* handler: tmp_uj__BrokerDatas type=*rocketmq.BrokerData kind=ptr quoted=false*/

				{
					if tok == fflib.FFTok_null {

						tmp_uj__BrokerDatas = nil

						state = fflib.FFParse_after_value
						goto mainparse
					}

					if tmp_uj__BrokerDatas == nil {
						tmp_uj__BrokerDatas = new(BrokerData)
					}

					err = tmp_uj__BrokerDatas.UnmarshalJSONFFLexer(fs, fflib.FFParse_want_key)
					if err != nil {
						return err
					}
					state = fflib.FFParse_after_value
				}

				uj.BrokerDatas = append(uj.BrokerDatas, tmp_uj__BrokerDatas)

				wantVal = false
			}
		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

wantedvalue:
	return fs.WrapErr(fmt.Errorf("wanted value token, but got token: %v", tok))
wrongtokenerror:
	return fs.WrapErr(fmt.Errorf("ffjson: wanted token: %v, but got token: %v output=%s", wantedTok, tok, fs.Output.String()))
tokerror:
	if fs.BigError != nil {
		return fs.WrapErr(fs.BigError)
	}
	err = fs.Error.ToError()
	if err != nil {
		return fs.WrapErr(err)
	}
	panic("ffjson-generated: unreachable, please report bug.")
done:

	return nil
}