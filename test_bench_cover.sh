#!/usr/bin/env bash
go test -cover -bench=. -benchtime=240s -cpuprofile=cpu.out -memprofile mem.out
go tool pprof rtree.test cpu.out
go tool pprof rtree.test mem.out

#type top20
