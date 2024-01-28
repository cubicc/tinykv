package server

import (
	"context"

	"github.com/pingcap-incubator/tinykv/kv/storage"
	"github.com/pingcap-incubator/tinykv/proto/pkg/kvrpcpb"
)

// The functions below are Server's Raw API. (implements TinyKvServer).
// Some helper methods can be found in sever.go in the current directory

// RawGet return the corresponding Get response based on RawGetRequest's CF and Key fields
func (server *Server) RawGet(_ context.Context, req *kvrpcpb.RawGetRequest) (*kvrpcpb.RawGetResponse, error) {
	// Your Code Here (1).
	r, err := server.storage.Reader(req.GetContext())
	if err != nil {
		return &kvrpcpb.RawGetResponse{Error: err.Error()}, err
	}

	val, err := r.GetCF(req.GetCf(), req.GetKey())
	if err != nil {
		return &kvrpcpb.RawGetResponse{Error: err.Error()}, err
	}

	resp := &kvrpcpb.RawGetResponse{
		Value:    val,
		NotFound: false,
	}
	if val == nil {
		resp.NotFound = true
	}
	return resp, nil
}

// RawPut puts the target data into storage and returns the corresponding response
func (server *Server) RawPut(_ context.Context, req *kvrpcpb.RawPutRequest) (*kvrpcpb.RawPutResponse, error) {
	// Your Code Here (1).
	// Hint: Consider using Storage.Modify to store data to be modified
	batch := []storage.Modify{{
		Data: storage.Put{
			Key:   req.GetKey(),
			Value: req.GetValue(),
			Cf:    req.GetCf(),
		},
	}}
	err := server.storage.Write(req.GetContext(), batch)
	if err != nil {
		return &kvrpcpb.RawPutResponse{Error: err.Error()}, err
	}

	return &kvrpcpb.RawPutResponse{}, nil
}

// RawDelete delete the target data from storage and returns the corresponding response
func (server *Server) RawDelete(_ context.Context, req *kvrpcpb.RawDeleteRequest) (*kvrpcpb.RawDeleteResponse, error) {
	// Your Code Here (1).
	// Hint: Consider using Storage.Modify to store data to be deleted
	batch := []storage.Modify{{
		Data: storage.Delete{
			Key: req.GetKey(),
			Cf:  req.GetCf(),
		},
	}}

	err := server.storage.Write(req.GetContext(), batch)
	if err != nil {
		return &kvrpcpb.RawDeleteResponse{Error: err.Error()}, err
	}
	return &kvrpcpb.RawDeleteResponse{}, nil
}

// RawScan scan the data starting from the start key up to limit. and return the corresponding result
func (server *Server) RawScan(_ context.Context, req *kvrpcpb.RawScanRequest) (*kvrpcpb.RawScanResponse, error) {
	// Your Code Here (1).
	// Hint: Consider using reader.IterCF
	if req.Limit == 0 {
		return &kvrpcpb.RawScanResponse{}, nil
	}

	r, err := server.storage.Reader(req.GetContext())
	if err != nil {
		return &kvrpcpb.RawScanResponse{Error: err.Error()}, err
	}
	iter := r.IterCF(req.GetCf())
	iter.Seek(req.StartKey)
	var pairs []*kvrpcpb.KvPair

	n := req.Limit
	for ; iter.Valid() && n > 0; iter.Next() {
		item := iter.Item()
		val, _ := item.ValueCopy(nil)
		pairs = append(pairs, &kvrpcpb.KvPair{
			Key:   item.KeyCopy(nil),
			Value: val,
		})
		n--
	}
	return &kvrpcpb.RawScanResponse{Kvs: pairs}, nil
}
