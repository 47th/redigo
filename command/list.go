package command

import (
	"redigo/resp"
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
		return HandleErr("ERR value is not a list", resp.SimpleError)
	}

	resEnv := resp.Envelope{
		OpCode:  resp.Integer,
		Integer: size,
	}

	return resp.FormatMapper(resEnv)

}

func rpush(env resp.Envelope, h *Handler) string {
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

	size, ok := h.store.RPUSH(key, values)

	if !ok {
		return HandleErr("ERR value is not a list", resp.SimpleError)
	}

	resEnv := resp.Envelope{
		OpCode:  resp.Integer,
		Integer: size,
	}

	return resp.FormatMapper(resEnv)

}

func lpop(env resp.Envelope, h *Handler) string {
	if env.Size != 2 {
		return HandleErr("ERR wrong number of arguments for 'incr' command", resp.SimpleError)
	}

	key := env.Array[1].String
	value, ok := h.store.LPOP(key)

	if !ok {
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
		return HandleErr("ERR wrong number of arguments for 'incr' command", resp.SimpleError)
	}

	key := env.Array[1].String
	value, ok := h.store.RPOP(key)

	if !ok {
		return HandleNIL()
	}

	resEnv := resp.Envelope{
		OpCode: resp.BulkString,
		String: value,
		Size:   len(value),
	}

	return resp.FormatMapper(resEnv)

}
