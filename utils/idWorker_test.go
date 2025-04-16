package utils

import (
	"sync"
	"testing"
	"time"
)
func TestGenerate(t *testing.T){
	worker,err:=NewIdWorker(1)
	if err!=nil{
		panic(err)
	}
	var wg sync.WaitGroup
	begin:=time.Now().UnixMilli()
	for range(300){
		wg.Add(1)
		go func(){
			for j:= 0 ;j < 100;j++{
				id,err:=worker.Generate()
				if err!=nil{
					panic(err)
				}
				println(id)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	end:=time.Now().UnixMilli()
	println("耗时：",end-begin,"ms")	

}