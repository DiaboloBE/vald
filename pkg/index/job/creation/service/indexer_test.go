// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package service

import (
	"context"
	"testing"

	agent "github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/pool"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/test/goleak"
	clientmock "github.com/vdaas/vald/internal/test/mock/client"
	grpcmock "github.com/vdaas/vald/internal/test/mock/grpc"
)

func Test_index_Start(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		client           discoverer.Client
		targetAddrs      []string
		creationPoolSize uint32
		concurrency      int
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			addrs := []string{
				"127.0.0.1:8080",
			}
			return test{
				name: "Success: when there is no error in the indexing request process",
				args: args{
					ctx: context.Background(),
				},
				fields: fields{
					client: &clientmock.DiscovererClientMock{
						GetAddrsFunc: func(_ context.Context) []string {
							return addrs
						},
						GetClientFunc: func() grpc.Client {
							return &grpcmock.GRPCClientMock{
								OrderedRangeConcurrentFunc: func(_ context.Context, _ []string, _ int,
									_ func(_ context.Context, _ string, _ *grpc.ClientConn, _ ...grpc.CallOption) error,
								) error {
									return nil
								},
							}
						},
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "Success: when a target addresses (targetAddrs) is given and there are no errors in the indexing request process",
				args: args{
					ctx: context.Background(),
				},
				fields: fields{
					targetAddrs: []string{
						"127.0.0.1:8080",
						"127.0.0.1:8081",
						"127.0.0.1:8083",
					},
					client: &clientmock.DiscovererClientMock{
						GetAddrsFunc: func(_ context.Context) []string {
							return nil
						},
						GetClientFunc: func() grpc.Client {
							return &grpcmock.GRPCClientMock{
								OrderedRangeConcurrentFunc: func(_ context.Context, _ []string, _ int,
									_ func(_ context.Context, _ string, _ *grpc.ClientConn, _ ...grpc.CallOption) error,
								) error {
									return nil
								},
								ConnectFunc: func(_ context.Context, _ string, _ ...grpc.DialOption) (pool.Conn, error) {
									return nil, nil
								},
							}
						},
					},
				},
			}
		}(),
		func() test {
			addrs := []string{
				"127.0.0.1:8080",
			}
			return test{
				name: "Fail: when there is an error wrapped with gRPC status in the indexing request process",
				args: args{
					ctx: context.Background(),
				},
				fields: fields{
					client: &clientmock.DiscovererClientMock{
						GetAddrsFunc: func(_ context.Context) []string {
							return addrs
						},
						GetClientFunc: func() grpc.Client {
							return &grpcmock.GRPCClientMock{
								OrderedRangeConcurrentFunc: func(_ context.Context, _ []string, _ int,
									_ func(_ context.Context, _ string, _ *grpc.ClientConn, _ ...grpc.CallOption) error,
								) error {
									return status.WrapWithInternal(
										agent.CreateIndexRPCName+" API connection not found",
										errors.ErrGRPCClientConnNotFound("*"),
									)
								},
							}
						},
					},
				},
				want: want{
					err: status.Error(codes.Internal,
						agent.CreateIndexRPCName+" API connection not found"),
				},
			}
		}(),
		func() test {
			addrs := []string{
				"127.0.0.1:8080",
			}
			return test{
				name: "Fail: When the OrderedRangeConcurrent method returns a gRPC client conn not found error",
				args: args{
					ctx: context.Background(),
				},
				fields: fields{
					client: &clientmock.DiscovererClientMock{
						GetAddrsFunc: func(_ context.Context) []string {
							return addrs
						},
						GetClientFunc: func() grpc.Client {
							return &grpcmock.GRPCClientMock{
								OrderedRangeConcurrentFunc: func(_ context.Context, _ []string, _ int,
									_ func(_ context.Context, _ string, _ *grpc.ClientConn, _ ...grpc.CallOption) error,
								) error {
									return errors.ErrGRPCClientConnNotFound("*")
								},
							}
						},
					},
				},
				want: want{
					err: status.Error(codes.Internal,
						agent.CreateIndexRPCName+" API connection not found"),
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			idx := &index{
				client:           test.fields.client,
				targetAddrs:      test.fields.targetAddrs,
				creationPoolSize: test.fields.creationPoolSize,
				concurrency:      test.fields.concurrency,
			}

			err := idx.Start(test.args.ctx)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		want Indexer
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Indexer, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Indexer, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           opts:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           opts:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got, err := New(test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_delDuplicateAddrs(t *testing.T) {
// 	type args struct {
// 		targetAddrs []string
// 	}
// 	type want struct {
// 		want []string
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, []string) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got []string) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           targetAddrs:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           targetAddrs:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := delDuplicateAddrs(test.args.targetAddrs)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_index_StartClient(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		client           discoverer.Client
// 		targetAddrs      []string
// 		creationPoolSize uint32
// 		concurrency      int
// 	}
// 	type want struct {
// 		want <-chan error
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, <-chan error, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got <-chan error, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           client:nil,
// 		           targetAddrs:nil,
// 		           creationPoolSize:0,
// 		           concurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           client:nil,
// 		           targetAddrs:nil,
// 		           creationPoolSize:0,
// 		           concurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			idx := &index{
// 				client:           test.fields.client,
// 				targetAddrs:      test.fields.targetAddrs,
// 				creationPoolSize: test.fields.creationPoolSize,
// 				concurrency:      test.fields.concurrency,
// 			}
//
// 			got, err := idx.StartClient(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_index_doCreateIndex(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		fn  func(_ context.Context, _ agent.AgentClient, _ ...grpc.CallOption) (*payload.Empty, error)
// 	}
// 	type fields struct {
// 		client           discoverer.Client
// 		targetAddrs      []string
// 		creationPoolSize uint32
// 		concurrency      int
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           fn:nil,
// 		       },
// 		       fields: fields {
// 		           client:nil,
// 		           targetAddrs:nil,
// 		           creationPoolSize:0,
// 		           concurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           fn:nil,
// 		           },
// 		           fields: fields {
// 		           client:nil,
// 		           targetAddrs:nil,
// 		           creationPoolSize:0,
// 		           concurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			idx := &index{
// 				client:           test.fields.client,
// 				targetAddrs:      test.fields.targetAddrs,
// 				creationPoolSize: test.fields.creationPoolSize,
// 				concurrency:      test.fields.concurrency,
// 			}
//
// 			err := idx.doCreateIndex(test.args.ctx, test.args.fn)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
