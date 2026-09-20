package command

import (
	"redigo/resp"
	"strconv"
)

func lpush(env resp.Envelope, h *Handler) string {
	if env.Size < 3 {
		return HandleErr("ERR wrong number of arguments for 'lpush' command", resp.SimpleError)
	}

	key := env.Array[1].String
	values := make([]string, env.Size-2)

	for i, e := range env.Array {
		if i > 1 {
			values[i-2] = e.String
		}
	}

	size, ok := h.store.LPUSH(key, values)

	if !ok {
		return HandleErr("WRONGTYPE Operation against a key holding the wrong kind of value", resp.SimpleError)
	}

	resEnv := resp.Envelope{
		OpCode:  resp.Integer,
		Integer: size,
	}

	return resp.FormatMapper(resEnv)

}

func rpush(env resp.Envelope, h *Handler) string {
	if env.Size < 3 {
		return HandleErr("ERR wrong number of arguments for 'rpush' command", resp.SimpleError)
	}

	key := env.Array[1].String
	values := make([]string, env.Size-2)

	for i, e := range env.Array {
		if i > 1 {
			values[i-2] = e.String
		}
	}

	size, ok := h.store.RPUSH(key, values)

	if !ok {
		return HandleErr("WRONGTYPE Operation against a key holding the wrong kind of value", resp.SimpleError)
	}

	resEnv := resp.Envelope{
		OpCode:  resp.Integer,
		Integer: size,
	}

	return resp.FormatMapper(resEnv)

}

func lpop(env resp.Envelope, h *Handler) string {
	if env.Size != 2 {
		return HandleErr("ERR wrong number of arguments for 'lpop' command", resp.SimpleError)
	}

	key := env.Array[1].String
	value, ok := h.store.LPOP(key)

	if !ok {
		if value == "WRONGTYPE" {
			return HandleErr("WRONGTYPE Operation against a key holding the wrong kind of value", resp.SimpleError)
		}
		return HandleNIL()
	}

	resEnv := resp.Envelope{
		OpCode: resp.BulkString,
		String: value,
		Size:   len(value),
	}

	return resp.FormatMapper(resEnv)

}

func rpop(env resp.Envelope, h *Handler) string {
	if env.Size != 2 {
		return HandleErr("ERR wrong number of arguments for 'rpop' command", resp.SimpleError)
	}

	key := env.Array[1].String
	value, ok := h.store.RPOP(key)

	if !ok {
		if value == "WRONGTYPE" {
			return HandleErr("WRONGTYPE Operation against a key holding the wrong kind of value", resp.SimpleError)
		}
		return HandleNIL()
	}

	resEnv := resp.Envelope{
		OpCode: resp.BulkString,
		String: value,
		Size:   len(value),
	}

	return resp.FormatMapper(resEnv)

}

func lrange(env resp.Envelope, h *Handler) string {
	if env.Size != 4 {
		return HandleErr("ERR wrong number of arguments for 'lrange' command", resp.SimpleError)
	}

	key := env.Array[1].String
	start, err := strconv.Atoi(env.Array[2].String)
	if err != nil {
		return HandleErr("ERR value is not an integer or out of range", resp.SimpleError)
	}
	stop, err := strconv.Atoi(env.Array[3].String)
	if err != nil {
		return HandleErr("ERR value is not an integer or out of range", resp.SimpleError)
	}

	values, ok := h.store.LRANGE(key, start, stop)
	if !ok {
		return HandleErr("WRONGTYPE Operation against a key holding the wrong kind of value", resp.SimpleError)
	}

	valEnv := make([]resp.Envelope, len(values))
	for i, val := range values {
		valEnv[i] = resp.Envelope{
			OpCode: resp.BulkString,
			Size:   len(val),
			String: val,
		}
	}

	resEnv := resp.Envelope{
		OpCode: resp.Array,
		Size:   len(values),
		Array:  valEnv,
	}

	return resp.FormatMapper(resEnv)
}

func llen(env resp.Envelope, h *Handler) string {
	if env.Size != 2 {
		return HandleErr("ERR wrong number of arguments for 'rpop' command", resp.SimpleError)
	}

	key := env.Array[1].String
	size, ok := h.store.LLEN(key)
	if !ok {
		return HandleErr("WRONGTYPE Operation against a key holding the wrong kind of value", resp.SimpleError)
	}

	resEnv := resp.Envelope{
		OpCode:  resp.Integer,
		Integer: size,
	}

	return resp.FormatMapper(resEnv)
}
