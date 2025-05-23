/*
 *     Copyright 2020 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package evaluator

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"

	commonv2 "d7y.io/api/v2/pkg/apis/common/v2"

	"d7y.io/dragonfly/v2/pkg/digest"
	"d7y.io/dragonfly/v2/pkg/idgen"
	"d7y.io/dragonfly/v2/pkg/types"
	"d7y.io/dragonfly/v2/scheduler/resource/standard"
)

var (
	mockRawHost = standard.Host{
		ID:              mockHostID,
		Type:            types.HostTypeNormal,
		Hostname:        "foo",
		IP:              "127.0.0.1",
		Port:            8003,
		DownloadPort:    8001,
		OS:              "darwin",
		Platform:        "darwin",
		PlatformFamily:  "Standalone Workstation",
		PlatformVersion: "11.1",
		KernelVersion:   "20.2.0",
		CPU:             mockCPU,
		Memory:          mockMemory,
		Network:         mockNetwork,
		Disk:            mockDisk,
		Build:           mockBuild,
		CreatedAt:       atomic.NewTime(time.Now()),
		UpdatedAt:       atomic.NewTime(time.Now()),
	}

	mockRawSeedHost = standard.Host{
		ID:              mockSeedHostID,
		Type:            types.HostTypeSuperSeed,
		Hostname:        "bar",
		IP:              "127.0.0.1",
		Port:            8003,
		DownloadPort:    8001,
		OS:              "darwin",
		Platform:        "darwin",
		PlatformFamily:  "Standalone Workstation",
		PlatformVersion: "11.1",
		KernelVersion:   "20.2.0",
		CPU:             mockCPU,
		Memory:          mockMemory,
		Network:         mockNetwork,
		Disk:            mockDisk,
		Build:           mockBuild,
		CreatedAt:       atomic.NewTime(time.Now()),
		UpdatedAt:       atomic.NewTime(time.Now()),
	}

	mockCPU = standard.CPU{
		LogicalCount:   4,
		PhysicalCount:  2,
		Percent:        1,
		ProcessPercent: 0.5,
		Times: standard.CPUTimes{
			User:      240662.2,
			System:    317950.1,
			Idle:      3393691.3,
			Nice:      0,
			Iowait:    0,
			Irq:       0,
			Softirq:   0,
			Steal:     0,
			Guest:     0,
			GuestNice: 0,
		},
	}

	mockMemory = standard.Memory{
		Total:              17179869184,
		Available:          5962813440,
		Used:               11217055744,
		UsedPercent:        65.291858,
		ProcessUsedPercent: 41.525125,
		Free:               2749598908,
	}

	mockNetwork = standard.Network{
		TCPConnectionCount:       10,
		UploadTCPConnectionCount: 1,
		Location:                 mockHostLocation,
		IDC:                      mockHostIDC,
		DownloadRate:             100,
		DownloadRateLimit:        200,
		UploadRate:               100,
		UploadRateLimit:          200,
	}

	mockDisk = standard.Disk{
		Total:             499963174912,
		Free:              37226479616,
		Used:              423809622016,
		UsedPercent:       91.92547406065952,
		InodesTotal:       4882452880,
		InodesUsed:        7835772,
		InodesFree:        4874617108,
		InodesUsedPercent: 0.1604884305611568,
		WriteBandwidth:    1,
		ReadBandwidth:     1,
	}

	mockBuild = standard.Build{
		GitVersion: "v1.0.0",
		GitCommit:  "221176b117c6d59366d68f2b34d38be50c935883",
		GoVersion:  "1.18",
		Platform:   "darwin",
	}

	mockTaskBackToSourceLimit   int32  = 200
	mockTaskURL                        = "http://example.com/foo"
	mockTaskPieceLength         uint64 = 2048
	mockTaskID                         = idgen.TaskIDV2ByURLBased(mockTaskURL, &mockTaskPieceLength, mockTaskTag, mockTaskApplication, mockTaskFilteredQueryParams)
	mockTaskDigest                     = digest.New(digest.AlgorithmSHA256, "c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4")
	mockTaskTag                        = "d7y"
	mockTaskApplication                = "foo"
	mockTaskFilteredQueryParams        = []string{"bar"}
	mockTaskHeader                     = map[string]string{"content-length": "100"}
	mockHostID                         = idgen.HostIDV2("127.0.0.1", "foo", false)
	mockSeedHostID                     = idgen.HostIDV2("127.0.0.1", "bar", true)
	mockHostLocation                   = "bas"
	mockHostIDC                        = "baz"
	mockPeerID                         = idgen.PeerIDV2()
)

func TestEvaluatorBase_newEvaluatorBase(t *testing.T) {
	tests := []struct {
		name   string
		expect func(t *testing.T, e any)
	}{
		{
			name: "new evaluator commonv1",
			expect: func(t *testing.T, e any) {
				assert := assert.New(t)
				assert.Equal(reflect.TypeOf(e).Elem().Name(), "evaluatorBase")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, newEvaluatorBase())
		})
	}
}

func TestEvaluatorBase_EvaluateParents(t *testing.T) {
	tests := []struct {
		name            string
		parents         []*standard.Peer
		child           *standard.Peer
		totalPieceCount uint32
		mock            func(parent []*standard.Peer, child *standard.Peer)
		expect          func(t *testing.T, parents []*standard.Peer)
	}{
		{
			name:    "parents is empty",
			parents: []*standard.Peer{},
			child: standard.NewPeer(idgen.PeerIDV1("127.0.0.1"),
				standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest)),
				standard.NewHost(
					mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
					mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.Type)),
			totalPieceCount: 1,
			mock: func(parent []*standard.Peer, child *standard.Peer) {
			},
			expect: func(t *testing.T, parents []*standard.Peer) {
				assert := assert.New(t)
				assert.Equal(len(parents), 0)

			},
		},
		{
			name: "evaluate single parent",
			parents: []*standard.Peer{
				standard.NewPeer(idgen.PeerIDV1("127.0.0.1"),
					standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest)),
					standard.NewHost(
						mockRawSeedHost.ID, mockRawSeedHost.IP, mockRawSeedHost.Hostname,
						mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.Type)),
			},
			child: standard.NewPeer(idgen.PeerIDV1("127.0.0.1"),
				standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest)),
				standard.NewHost(
					mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
					mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.Type)),
			totalPieceCount: 1,
			mock: func(parent []*standard.Peer, child *standard.Peer) {
			},
			expect: func(t *testing.T, parents []*standard.Peer) {
				assert := assert.New(t)
				assert.Equal(len(parents), 1)
				assert.Equal(parents[0].Task.ID, mockTaskID)
				assert.Equal(parents[0].Host.ID, mockRawSeedHost.ID)

			},
		},
		{
			name: "evaluate parents with free upload count",
			parents: []*standard.Peer{
				standard.NewPeer(idgen.PeerIDV1("127.0.0.1"),
					standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest)),
					standard.NewHost(
						mockRawSeedHost.ID, mockRawSeedHost.IP, mockRawSeedHost.Hostname,
						mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.Type)),
				standard.NewPeer(idgen.PeerIDV1("127.0.0.1"),
					standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest)),
					standard.NewHost(
						"bar", mockRawSeedHost.IP, mockRawSeedHost.Hostname,
						mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.Type)),
				standard.NewPeer(idgen.PeerIDV1("127.0.0.1"),
					standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest)),
					standard.NewHost(
						"baz", mockRawSeedHost.IP, mockRawSeedHost.Hostname,
						mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.Type)),
				standard.NewPeer(idgen.PeerIDV1("127.0.0.1"),
					standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest)),
					standard.NewHost(
						"bac", mockRawSeedHost.IP, mockRawSeedHost.Hostname,
						mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.Type)),
				standard.NewPeer(idgen.PeerIDV1("127.0.0.1"),
					standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest)),
					standard.NewHost(
						"bae", mockRawSeedHost.IP, mockRawSeedHost.Hostname,
						mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.Type)),
			},
			child: standard.NewPeer(idgen.PeerIDV1("127.0.0.1"),
				standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest)),
				standard.NewHost(
					mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
					mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.Type)),
			totalPieceCount: 1,
			mock: func(parents []*standard.Peer, child *standard.Peer) {
				parents[1].Host.ConcurrentUploadCount.Add(4)
				parents[2].Host.ConcurrentUploadCount.Add(3)
				parents[3].Host.ConcurrentUploadCount.Add(2)
				parents[4].Host.ConcurrentUploadCount.Add(1)
			},
			expect: func(t *testing.T, parents []*standard.Peer) {
				assert := assert.New(t)
				assert.Equal(len(parents), 5)
				assert.Equal(parents[0].Host.ID, mockRawSeedHost.ID)
				assert.Equal(parents[1].Host.ID, "bae")
				assert.Equal(parents[2].Host.ID, "bac")
				assert.Equal(parents[3].Host.ID, "baz")
				assert.Equal(parents[4].Host.ID, "bar")
			},
		},
		{
			name: "evaluate parents with pieces",
			parents: []*standard.Peer{
				standard.NewPeer(idgen.PeerIDV1("127.0.0.1"),
					standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest)),
					standard.NewHost(
						mockRawSeedHost.ID, mockRawSeedHost.IP, mockRawSeedHost.Hostname,
						mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.Type)),
				standard.NewPeer(idgen.PeerIDV1("127.0.0.1"),
					standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest)),
					standard.NewHost(
						"bar", mockRawSeedHost.IP, mockRawSeedHost.Hostname,
						mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.Type)),
				standard.NewPeer(idgen.PeerIDV1("127.0.0.1"),
					standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest)),
					standard.NewHost(
						"baz", mockRawSeedHost.IP, mockRawSeedHost.Hostname,
						mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.Type)),
				standard.NewPeer(idgen.PeerIDV1("127.0.0.1"),
					standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest)),
					standard.NewHost(
						"bac", mockRawSeedHost.IP, mockRawSeedHost.Hostname,
						mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.Type)),
				standard.NewPeer(idgen.PeerIDV1("127.0.0.1"),
					standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest)),
					standard.NewHost(
						"bae", mockRawSeedHost.IP, mockRawSeedHost.Hostname,
						mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.Type)),
			},
			child: standard.NewPeer(idgen.PeerIDV1("127.0.0.1"),
				standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest)),
				standard.NewHost(
					mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
					mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.Type)),
			totalPieceCount: 1,
			mock: func(parents []*standard.Peer, child *standard.Peer) {
				parents[1].FinishedPieces.Set(0)
				parents[2].FinishedPieces.Set(0).Set(1)
				parents[3].FinishedPieces.Set(0).Set(1).Set(2)
				parents[4].FinishedPieces.Set(0).Set(1).Set(2).Set(3)
			},
			expect: func(t *testing.T, parents []*standard.Peer) {
				assert := assert.New(t)
				assert.Equal(len(parents), 5)
				assert.Equal(parents[0].Host.ID, "bae")
				assert.Equal(parents[1].Host.ID, "bac")
				assert.Equal(parents[2].Host.ID, "baz")
				assert.Equal(parents[3].Host.ID, "bar")
				assert.Equal(parents[4].Host.ID, mockRawSeedHost.ID)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := newEvaluatorBase()
			tc.mock(tc.parents, tc.child)
			tc.expect(t, e.EvaluateParents(tc.parents, tc.child, tc.totalPieceCount))
		})
	}
}

func TestEvaluatorBase_evaluate(t *testing.T) {
	tests := []struct {
		name            string
		parent          *standard.Peer
		child           *standard.Peer
		totalPieceCount uint32
		mock            func(parent *standard.Peer, child *standard.Peer)
		expect          func(t *testing.T, score float64)
	}{
		{
			name: "evaluate parent",
			parent: standard.NewPeer(idgen.PeerIDV1("127.0.0.1"),
				standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest)),
				standard.NewHost(
					mockRawSeedHost.ID, mockRawSeedHost.IP, mockRawSeedHost.Hostname,
					mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.Type)),
			child: standard.NewPeer(idgen.PeerIDV1("127.0.0.1"),
				standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest)),
				standard.NewHost(
					mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
					mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.Type)),
			totalPieceCount: 1,
			mock: func(parent *standard.Peer, child *standard.Peer) {
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0.35))
			},
		},
		{
			name: "evaluate parent with pieces",
			parent: standard.NewPeer(idgen.PeerIDV1("127.0.0.1"),
				standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest)),
				standard.NewHost(
					mockRawSeedHost.ID, mockRawSeedHost.IP, mockRawSeedHost.Hostname,
					mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.Type)),
			child: standard.NewPeer(idgen.PeerIDV1("127.0.0.1"),
				standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest)),
				standard.NewHost(
					mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
					mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.Type)),
			totalPieceCount: 1,
			mock: func(parent *standard.Peer, child *standard.Peer) {
				parent.FinishedPieces.Set(0)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0.55))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := newEvaluatorBase()
			tc.mock(tc.parent, tc.child)
			tc.expect(t, e.(*evaluatorBase).evaluateParents(tc.parent, tc.child, tc.totalPieceCount))
		})
	}
}

func TestEvaluatorBase_calculatePieceScore(t *testing.T) {
	mockHost := standard.NewHost(
		mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
		mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.Type)
	mockTask := standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest))

	tests := []struct {
		name            string
		parent          *standard.Peer
		child           *standard.Peer
		totalPieceCount uint32
		mock            func(parent *standard.Peer, child *standard.Peer)
		expect          func(t *testing.T, score float64)
	}{
		{
			name:            "total piece count is zero and child pieces are empty",
			parent:          standard.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			child:           standard.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			totalPieceCount: 0,
			mock: func(parent *standard.Peer, child *standard.Peer) {
				parent.FinishedPieces.Set(0)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(1))
			},
		},
		{
			name:            "total piece count is zero and parent pieces are empty",
			parent:          standard.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			child:           standard.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			totalPieceCount: 0,
			mock: func(parent *standard.Peer, child *standard.Peer) {
				child.FinishedPieces.Set(0)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(-1))
			},
		},
		{
			name:            "total piece count is zero and child pieces of length greater than parent pieces",
			parent:          standard.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			child:           standard.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			totalPieceCount: 0,
			mock: func(parent *standard.Peer, child *standard.Peer) {
				parent.FinishedPieces.Set(0)
				child.FinishedPieces.Set(0)
				child.FinishedPieces.Set(1)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(-1))
			},
		},
		{
			name:            "total piece count is zero and child pieces of length equal than parent pieces",
			parent:          standard.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			child:           standard.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			totalPieceCount: 0,
			mock: func(parent *standard.Peer, child *standard.Peer) {
				parent.FinishedPieces.Set(0)
				child.FinishedPieces.Set(0)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name:            "total piece count is zero and parent pieces of length greater than child pieces",
			parent:          standard.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			child:           standard.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			totalPieceCount: 0,
			mock: func(parent *standard.Peer, child *standard.Peer) {
				parent.FinishedPieces.Set(0)
				parent.FinishedPieces.Set(1)
				child.FinishedPieces.Set(0)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(1))
			},
		},
		{
			name:            "parent pieces are empty",
			parent:          standard.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			child:           standard.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			totalPieceCount: 10,
			mock:            func(parent *standard.Peer, child *standard.Peer) {},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name:            "parent pieces of length greater than zero",
			parent:          standard.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			child:           standard.NewPeer(idgen.PeerIDV1("127.0.0.1"), mockTask, mockHost),
			totalPieceCount: 10,
			mock: func(parent *standard.Peer, child *standard.Peer) {
				parent.FinishedPieces.Set(0)
				parent.FinishedPieces.Set(1)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0.2))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := newEvaluatorBase()
			tc.mock(tc.parent, tc.child)
			tc.expect(t, e.(*evaluatorBase).calculatePieceScore(tc.parent.FinishedPieces.Count(), tc.child.FinishedPieces.Count(), tc.totalPieceCount))
		})
	}
}

func TestEvaluatorBase_calculatehostUploadSuccessScore(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(host *standard.Host)
		expect func(t *testing.T, score float64)
	}{
		{
			name: "UploadFailedCount is larger than UploadCount",
			mock: func(host *standard.Host) {
				host.UploadCount.Add(1)
				host.UploadFailedCount.Add(2)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "UploadFailedCount and UploadCount is zero",
			mock: func(host *standard.Host) {
				host.UploadCount.Add(0)
				host.UploadFailedCount.Add(0)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(1))
			},
		},
		{
			name: "UploadCount is larger than UploadFailedCount",
			mock: func(host *standard.Host) {
				host.UploadCount.Add(2)
				host.UploadFailedCount.Add(1)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0.5))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			host := standard.NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.Type)
			mockTask := standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest))
			mockPeer := standard.NewPeer(mockPeerID, mockTask, host)
			e := newEvaluatorBase()
			tc.mock(host)
			tc.expect(t, e.(*evaluatorBase).calculateParentHostUploadSuccessScore(mockPeer.Host.UploadCount.Load(), mockPeer.Host.UploadFailedCount.Load()))
		})
	}
}

func TestEvaluatorBase_calculateFreeUploadScore(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(host *standard.Host, mockPeer *standard.Peer)
		expect func(t *testing.T, score float64)
	}{
		{
			name: "host peers is not empty",
			mock: func(host *standard.Host, mockPeer *standard.Peer) {
				mockPeer.Host.ConcurrentUploadCount.Add(1)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0.995))
			},
		},
		{
			name: "host peers is empty",
			mock: func(host *standard.Host, mockPeer *standard.Peer) {},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(1))
			},
		},
		{
			name: "freeUploadCount is empty",
			mock: func(host *standard.Host, mockPeer *standard.Peer) {
				mockPeer.Host.ConcurrentUploadCount.Add(host.ConcurrentUploadLimit.Load())
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			host := standard.NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.Type)
			mockTask := standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest))
			mockPeer := standard.NewPeer(mockPeerID, mockTask, host)
			e := newEvaluatorBase()
			tc.mock(host, mockPeer)
			tc.expect(t, e.(*evaluatorBase).calculateFreeUploadScore(host))
		})
	}
}

func TestEvaluatorBase_calculateHostTypeScore(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(peer *standard.Peer)
		expect func(t *testing.T, score float64)
	}{
		{
			name: "host is normal peer",
			mock: func(peer *standard.Peer) {},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0.5))
			},
		},
		{
			name: "host is seed peer but peer state is PeerStateSucceeded",
			mock: func(peer *standard.Peer) {
				peer.FSM.SetState(standard.PeerStateSucceeded)
				peer.Host.Type = types.HostTypeSuperSeed
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "host is seed peer but peer state is PeerStateRunning",
			mock: func(peer *standard.Peer) {
				peer.FSM.SetState(standard.PeerStateRunning)
				peer.Host.Type = types.HostTypeSuperSeed
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(1))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockHost := standard.NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.Type)
			mockTask := standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest))
			peer := standard.NewPeer(mockPeerID, mockTask, mockHost)
			e := newEvaluatorBase()
			tc.mock(peer)
			tc.expect(t, e.(*evaluatorBase).calculateHostTypeScore(peer))
		})
	}
}

func TestEvaluatorBase_calculateIDCAffinityScore(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(dstHost *standard.Host, srcHost *standard.Host)
		expect func(t *testing.T, score float64)
	}{
		{
			name: "idc is empty",
			mock: func(dstHost *standard.Host, srcHost *standard.Host) {
				dstHost.Network.IDC = ""
				srcHost.Network.IDC = ""
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "dst host idc is empty",
			mock: func(dstHost *standard.Host, srcHost *standard.Host) {
				dstHost.Network.IDC = ""
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "src host idc is empty",
			mock: func(dstHost *standard.Host, srcHost *standard.Host) {
				srcHost.Network.IDC = ""
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "idc is not the same",
			mock: func(dstHost *standard.Host, srcHost *standard.Host) {
				dstHost.Network.IDC = "foo"
				srcHost.Network.IDC = "bar"
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "idc is the same",
			mock: func(dstHost *standard.Host, srcHost *standard.Host) {
				dstHost.Network.IDC = "example"
				srcHost.Network.IDC = "example"
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(1))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dstHost := standard.NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.Type)
			srcHost := standard.NewHost(
				mockRawSeedHost.ID, mockRawSeedHost.IP, mockRawSeedHost.Hostname,
				mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.Type)
			e := newEvaluatorBase()
			tc.mock(dstHost, srcHost)
			tc.expect(t, e.(*evaluatorBase).calculateIDCAffinityScore(dstHost.Network.IDC, srcHost.Network.IDC))
		})
	}
}

func TestEvaluatorBase_calculateMultiElementAffinityScore(t *testing.T) {
	tests := []struct {
		name   string
		dst    string
		src    string
		expect func(t *testing.T, score float64)
	}{
		{
			name: "dst is empty and src is empty",
			dst:  "",
			src:  "",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "dst is empty",
			dst:  "",
			src:  "baz",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "src is empty",
			dst:  "bar",
			src:  "",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "has only one element and matches",
			dst:  "foo",
			src:  "foo",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(1))
			},
		},
		{
			name: "has only one element and does not match",
			dst:  "foo",
			src:  "bar",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "has multi element and match",
			dst:  "foo|bar",
			src:  "foo|bar",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(1))
			},
		},
		{
			name: "has multi element and does not match",
			dst:  "foo|bar",
			src:  "bar|foo",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "dst length is greater than src",
			dst:  "foo|bar|baz",
			src:  "foo|bar",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0.4))
			},
		},
		{
			name: "src length is greater than dst",
			dst:  "foo|bar",
			src:  "foo|bar|baz",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0.4))
			},
		},
		{
			name: "dst exceeds maximum length",
			dst:  "foo|bar|baz|bac|bae|baf",
			src:  "foo|bar|baz",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0.6))
			},
		},
		{
			name: "src exceeds maximum length",
			dst:  "foo|bar|baz",
			src:  "foo|bar|baz|bac|bae|baf",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0.6))
			},
		},
		{
			name: "dst and src both exceeds maximum length",
			dst:  "foo|bar|baz|bac|bae|baf",
			src:  "foo|bar|baz|bac|bae|bai",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(1))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := newEvaluatorBase()
			tc.expect(t, e.(*evaluatorBase).calculateMultiElementAffinityScore(tc.dst, tc.src))
		})
	}
}

func TestEvaluatorBase_IsBadParent(t *testing.T) {
	mockHost := standard.NewHost(
		mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
		mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.Type)
	mockTask := standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest))

	tests := []struct {
		name            string
		peer            *standard.Peer
		totalPieceCount int32
		mock            func(peer *standard.Peer)
		expect          func(t *testing.T, isBadParent bool)
	}{
		{
			name:            "peer state is PeerStateFailed",
			peer:            standard.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *standard.Peer) {
				peer.FSM.SetState(standard.PeerStateFailed)
			},
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.True(isBadParent)
			},
		},
		{
			name:            "peer state is PeerStateLeave",
			peer:            standard.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *standard.Peer) {
				peer.FSM.SetState(standard.PeerStateLeave)
			},
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.True(isBadParent)
			},
		},
		{
			name:            "peer state is PeerStatePending",
			peer:            standard.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *standard.Peer) {
				peer.FSM.SetState(standard.PeerStatePending)
			},
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.True(isBadParent)
			},
		},
		{
			name:            "peer state is PeerStateReceivedTiny",
			peer:            standard.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *standard.Peer) {
				peer.FSM.SetState(standard.PeerStateReceivedTiny)
			},
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.True(isBadParent)
			},
		},
		{
			name:            "peer state is PeerStateReceivedSmall",
			peer:            standard.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *standard.Peer) {
				peer.FSM.SetState(standard.PeerStateReceivedSmall)
			},
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.True(isBadParent)
			},
		},
		{
			name:            "peer state is PeerStateReceivedNormal",
			peer:            standard.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *standard.Peer) {
				peer.FSM.SetState(standard.PeerStateReceivedNormal)
			},
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.True(isBadParent)
			},
		},
		{
			name:            "download costs does not meet the normal distribution and last cost is twenty times more than mean",
			peer:            standard.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *standard.Peer) {
				peer.FSM.SetState(standard.PeerStateRunning)
				peer.AppendPieceCost(10)
				peer.AppendPieceCost(201)
			},
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.True(isBadParent)
			},
		},
		{
			name:            "download costs does not meet the normal distribution and last cost is twenty times lower than mean",
			peer:            standard.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *standard.Peer) {
				peer.FSM.SetState(standard.PeerStateRunning)
				peer.AppendPieceCost(10)
				peer.AppendPieceCost(200)
			},
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.False(isBadParent)
			},
		},
		{
			name:            "download costs meet the normal distribution and last cost is too long",
			peer:            standard.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *standard.Peer) {
				peer.FSM.SetState(standard.PeerStateRunning)
				for i := range 30 {
					peer.AppendPieceCost(time.Duration(i))
				}
				peer.AppendPieceCost(50)
			},
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.True(isBadParent)
			},
		},
		{
			name:            "download costs meet the normal distribution and last cost is normal",
			peer:            standard.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *standard.Peer) {
				peer.FSM.SetState(standard.PeerStateRunning)
				for i := range 30 {
					peer.AppendPieceCost(time.Duration(i))
				}
				peer.AppendPieceCost(18)
			},
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.False(isBadParent)
			},
		},
		{
			name:            "download costs meet the normal distribution and last cost is too short",
			peer:            standard.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *standard.Peer) {
				peer.FSM.SetState(standard.PeerStateRunning)
				for i := 20; i < 50; i++ {
					peer.AppendPieceCost(time.Duration(i))
				}
				peer.AppendPieceCost(0)
			},
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.False(isBadParent)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := newEvaluatorBase()
			tc.mock(tc.peer)
			tc.expect(t, e.IsBadParent(tc.peer))
		})
	}
}
