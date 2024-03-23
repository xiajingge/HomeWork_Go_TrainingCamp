package main

import (
	"fmt"
)

/*
*
实现删除切片特定下标元素的方法。
要求一：能够实现删除操作就可以。
要求二：考虑使用比较高性能的实现。
要求三：改造为泛型方法
要求四：支持缩容，并旦设计缩容机制
*/
func Delete[T any](src []T, index int) ([]T, T, error) {
	length := len(src)
	// 判断要删除的索引是否合理
	if index >= length || index < 0 {
		var zero T
		return nil, zero, newErrIndexOutOfRange(length, index)
	}
	//将要删除索引的后面的值向前挪一位
	res := src[index]
	for i := index; i+1 < length; i++ {
		src[i] = src[i+1]
	}
	//因为没有新开辟一个空间，因此挪移后，长度不变，所以要将最后一位的数据删除
	src = src[:length-1]
	src = Shrink(src)
	return src, res, nil
}

/*
*
根据特定的阈值和比例因子来动态调整给定的容量值
*/
func calCapacity(c, l int) (int, bool) {
	// 当容量小于64时，不用缩容
	if c <= 64 {
		return c, false
	}
	// 当容量大于2048并且容量是长度的两倍以上时，根据比例因子缩容
	if c > 2048 && (c/l >= 2) {
		factor := 0.625
		return int(float32(c) * float32(factor)), true
	}
	// 当容量小于2048并且容量是长度的四倍以上时，缩容为原来的一半容量
	if c <= 2048 && (c/l >= 4) {
		return c / 2, true
	}
	return c, false
}

/*
*
将给定的切片进行缩容
*/
func Shrink[T any](src []T) []T {
	c, l := cap(src), len(src)
	n, changed := calCapacity(c, l)
	if !changed {
		s := make([]T, 0, n)
		s = append(s, src...)
		return s
	}
	return src
}

func newErrIndexOutOfRange(length, index int) error {
	return fmt.Errorf("index %d out of length %d", index, length)
}
