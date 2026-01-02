package ocppj_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/xBlaz3kx/ocpp-go/ocppj"
)

// BenchmarkFIFOClientQueue_Push benchmarks pushing elements to the queue
func BenchmarkFIFOClientQueue_Push(b *testing.B) {
	queue := ocppj.NewFIFOClientQueue(0) // Unlimited capacity
	req := newMockRequest("benchmark")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = queue.Push(req)
	}
}

// BenchmarkFIFOClientQueue_PushParallel benchmarks concurrent pushes
func BenchmarkFIFOClientQueue_PushParallel(b *testing.B) {
	queue := ocppj.NewFIFOClientQueue(0) // Unlimited capacity
	req := newMockRequest("benchmark")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = queue.Push(req)
		}
	})
}

// BenchmarkFIFOClientQueue_Pop benchmarks popping elements from the queue
func BenchmarkFIFOClientQueue_Pop(b *testing.B) {
	queue := ocppj.NewFIFOClientQueue(0)
	req := newMockRequest("benchmark")

	// Pre-populate queue
	for i := 0; i < b.N; i++ {
		_ = queue.Push(req)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = queue.Pop()
	}
}

// BenchmarkFIFOClientQueue_PopParallel benchmarks concurrent pops
func BenchmarkFIFOClientQueue_PopParallel(b *testing.B) {
	queue := ocppj.NewFIFOClientQueue(0)
	req := newMockRequest("benchmark")

	// Pre-populate queue
	for i := 0; i < b.N; i++ {
		_ = queue.Push(req)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = queue.Pop()
		}
	})
}

// BenchmarkFIFOClientQueue_Peek benchmarks peeking at the front element
func BenchmarkFIFOClientQueue_Peek(b *testing.B) {
	queue := ocppj.NewFIFOClientQueue(0)
	req := newMockRequest("benchmark")
	_ = queue.Push(req)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = queue.Peek()
	}
}

// BenchmarkFIFOClientQueue_PeekParallel benchmarks concurrent peeks
func BenchmarkFIFOClientQueue_PeekParallel(b *testing.B) {
	queue := ocppj.NewFIFOClientQueue(0)
	req := newMockRequest("benchmark")
	_ = queue.Push(req)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = queue.Peek()
		}
	})
}

// BenchmarkFIFOClientQueue_Size benchmarks getting the queue size
func BenchmarkFIFOClientQueue_Size(b *testing.B) {
	queue := ocppj.NewFIFOClientQueue(0)
	req := newMockRequest("benchmark")
	for i := 0; i < 1000; i++ {
		_ = queue.Push(req)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = queue.Size()
	}
}

// BenchmarkFIFOClientQueue_SizeParallel benchmarks concurrent size checks
func BenchmarkFIFOClientQueue_SizeParallel(b *testing.B) {
	queue := ocppj.NewFIFOClientQueue(0)
	req := newMockRequest("benchmark")
	for i := 0; i < 1000; i++ {
		_ = queue.Push(req)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = queue.Size()
		}
	})
}

// BenchmarkFIFOClientQueue_IsEmpty benchmarks checking if queue is empty
func BenchmarkFIFOClientQueue_IsEmpty(b *testing.B) {
	queue := ocppj.NewFIFOClientQueue(0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = queue.IsEmpty()
	}
}

// BenchmarkFIFOClientQueue_IsFull benchmarks checking if queue is full
func BenchmarkFIFOClientQueue_IsFull(b *testing.B) {
	capacity := 1000
	queue := ocppj.NewFIFOClientQueue(capacity)
	req := newMockRequest("benchmark")
	for i := 0; i < capacity; i++ {
		_ = queue.Push(req)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = queue.IsFull()
	}
}

// BenchmarkFIFOClientQueue_PushPop benchmarks the push-pop cycle
func BenchmarkFIFOClientQueue_PushPop(b *testing.B) {
	queue := ocppj.NewFIFOClientQueue(0)
	req := newMockRequest("benchmark")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = queue.Push(req)
		_ = queue.Pop()
	}
}

// BenchmarkFIFOClientQueue_PushPopParallel benchmarks concurrent push-pop operations
func BenchmarkFIFOClientQueue_PushPopParallel(b *testing.B) {
	queue := ocppj.NewFIFOClientQueue(0)
	req := newMockRequest("benchmark")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = queue.Push(req)
			_ = queue.Pop()
		}
	})
}

// BenchmarkFIFOClientQueue_MixedOperations benchmarks mixed operations
func BenchmarkFIFOClientQueue_MixedOperations(b *testing.B) {
	queue := ocppj.NewFIFOClientQueue(1000)
	req := newMockRequest("benchmark")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = queue.Push(req)
		_ = queue.Peek()
		_ = queue.Size()
		_ = queue.IsEmpty()
		if i%2 == 0 {
			_ = queue.Pop()
		}
	}
}

// BenchmarkFIFOQueueMap_Get benchmarks getting a queue from the map
func BenchmarkFIFOQueueMap_Get(b *testing.B) {
	queueMap := ocppj.NewFIFOQueueMap(100)
	clientID := "client1"
	queueMap.GetOrCreate(clientID)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = queueMap.Get(clientID)
	}
}

// BenchmarkFIFOQueueMap_GetParallel benchmarks concurrent gets
func BenchmarkFIFOQueueMap_GetParallel(b *testing.B) {
	queueMap := ocppj.NewFIFOQueueMap(100)
	clientID := "client1"
	queueMap.GetOrCreate(clientID)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = queueMap.Get(clientID)
		}
	})
}

// BenchmarkFIFOQueueMap_GetOrCreate benchmarks getting or creating a queue
func BenchmarkFIFOQueueMap_GetOrCreate(b *testing.B) {
	queueMap := ocppj.NewFIFOQueueMap(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		clientID := fmt.Sprintf("client%d", i%10) // Reuse 10 clients
		_ = queueMap.GetOrCreate(clientID)
	}
}

// BenchmarkFIFOQueueMap_GetOrCreateParallel benchmarks concurrent get-or-create operations
func BenchmarkFIFOQueueMap_GetOrCreateParallel(b *testing.B) {
	queueMap := ocppj.NewFIFOQueueMap(100)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			clientID := fmt.Sprintf("client%d", i%10) // Reuse 10 clients
			_ = queueMap.GetOrCreate(clientID)
			i++
		}
	})
}

// BenchmarkFIFOQueueMap_Add benchmarks adding a queue to the map
func BenchmarkFIFOQueueMap_Add(b *testing.B) {
	queueMap := ocppj.NewFIFOQueueMap(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		clientID := fmt.Sprintf("client%d", i)
		queue := ocppj.NewFIFOClientQueue(100)
		queueMap.Add(clientID, queue)
	}
}

// BenchmarkFIFOQueueMap_Remove benchmarks removing a queue from the map
func BenchmarkFIFOQueueMap_Remove(b *testing.B) {
	queueMap := ocppj.NewFIFOQueueMap(100)

	// Pre-populate map
	for i := 0; i < b.N; i++ {
		clientID := fmt.Sprintf("client%d", i)
		queueMap.GetOrCreate(clientID)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		clientID := fmt.Sprintf("client%d", i)
		queueMap.Remove(clientID)
	}
}

// BenchmarkFIFOQueueMap_Size benchmarks getting total size across all queues
func BenchmarkFIFOQueueMap_Size(b *testing.B) {
	queueMap := ocppj.NewFIFOQueueMap(100)
	req := newMockRequest("benchmark")

	// Pre-populate with queues and elements
	for i := 0; i < 100; i++ {
		clientID := fmt.Sprintf("client%d", i)
		queue := queueMap.GetOrCreate(clientID)
		for j := 0; j < 10; j++ {
			_ = queue.Push(req)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = queueMap.Size()
	}
}

// BenchmarkFIFOQueueMap_SizeParallel benchmarks concurrent size checks
func BenchmarkFIFOQueueMap_SizeParallel(b *testing.B) {
	queueMap := ocppj.NewFIFOQueueMap(100)
	req := newMockRequest("benchmark")

	// Pre-populate with queues and elements
	for i := 0; i < 100; i++ {
		clientID := fmt.Sprintf("client%d", i)
		queue := queueMap.GetOrCreate(clientID)
		for j := 0; j < 10; j++ {
			_ = queue.Push(req)
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = queueMap.Size()
		}
	})
}

// BenchmarkFIFOQueueMap_SizePerClient benchmarks getting size per client
func BenchmarkFIFOQueueMap_SizePerClient(b *testing.B) {
	queueMap := ocppj.NewFIFOQueueMap(100)
	req := newMockRequest("benchmark")

	// Pre-populate with queues and elements
	for i := 0; i < 100; i++ {
		clientID := fmt.Sprintf("client%d", i)
		queue := queueMap.GetOrCreate(clientID)
		for j := 0; j < 10; j++ {
			_ = queue.Push(req)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = queueMap.SizePerClient()
	}
}

// BenchmarkFIFOQueueMap_SizePerClientParallel benchmarks concurrent size-per-client checks
func BenchmarkFIFOQueueMap_SizePerClientParallel(b *testing.B) {
	queueMap := ocppj.NewFIFOQueueMap(100)
	req := newMockRequest("benchmark")

	// Pre-populate with queues and elements
	for i := 0; i < 100; i++ {
		clientID := fmt.Sprintf("client%d", i)
		queue := queueMap.GetOrCreate(clientID)
		for j := 0; j < 10; j++ {
			_ = queue.Push(req)
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = queueMap.SizePerClient()
		}
	})
}

// BenchmarkFIFOQueueMap_MixedOperations benchmarks mixed map operations
func BenchmarkFIFOQueueMap_MixedOperations(b *testing.B) {
	queueMap := ocppj.NewFIFOQueueMap(100)
	req := newMockRequest("benchmark")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		clientID := fmt.Sprintf("client%d", i%10)
		queue := queueMap.GetOrCreate(clientID)
		_ = queue.Push(req)
		_ = queueMap.Size()
		if i%5 == 0 {
			queueMap.Remove(clientID)
		}
	}
}

// BenchmarkFIFOQueueMap_ConcurrentClients benchmarks concurrent operations on different clients
func BenchmarkFIFOQueueMap_ConcurrentClients(b *testing.B) {
	queueMap := ocppj.NewFIFOQueueMap(100)
	req := newMockRequest("benchmark")
	numClients := 10

	// Pre-create clients
	for i := 0; i < numClients; i++ {
		clientID := fmt.Sprintf("client%d", i)
		queueMap.GetOrCreate(clientID)
	}

	b.ResetTimer()
	var wg sync.WaitGroup
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(clientID string) {
			defer wg.Done()
			for j := 0; j < b.N/numClients; j++ {
				queue, _ := queueMap.Get(clientID)
				if queue != nil {
					_ = queue.Push(req)
					_ = queue.Pop()
				}
			}
		}(fmt.Sprintf("client%d", i))
	}
	wg.Wait()
}
