package rpcc

import (
	"io"

	"google.golang.org/grpc"

	"github.com/fidelfly/gostool/grpcx/protob"
)

func SendChunk(stream grpc.ClientStream, reader io.Reader, chunkSize ...int) error {
	size := 1024
	if len(chunkSize) > 0 {
		size = chunkSize[0]
	}

	buf := make([]byte, size)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		stream.SendMsg(&protob.Chunk{
			Chunk: buf[:n],
		})
	}

	return nil
}

func ReceiveChunk(stream grpc.ClientStream, writer io.Writer) error {
	for {
		c := new(protob.Chunk)
		if err := stream.RecvMsg(c); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		writer.Write(c.Chunk)
	}

	return nil
}
