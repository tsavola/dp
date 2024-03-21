// Copyright (c) 2024 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package sync

import (
	"sync"
)

type Map[K comparable, V any] struct {
	m sync.Map
}

func (m *Map[K, V]) LoadOrStore(k K, v V) V {
	x, _ := m.m.LoadOrStore(k, v)
	return x.(V)
}
