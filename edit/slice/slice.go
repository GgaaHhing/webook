package slice

import "errors"

// DeleteSli 删除
func DeleteSli[T any](src []T, index int) ([]T, error) {
	if index < 0 || index >= len(src) {
		return nil, errors.New("index 超出 切片范围")
	}
	new_src := make([]T, len(src)-1)
	//append 会改变原切片所引用的数组的内容（如果容量不足，会创建新的底层数组）。
	//copy 不会改变原切片的内容，只是将原切片的内容复制到另一个切片中。
	copy(new_src, src[:index])
	copy(new_src[index:], src[index+1:])
	return Shrink(new_src), nil
}

// Shrink 缩容
func Shrink[T any](src []T) []T {
	// len: 切片的长度表示切片中当前包含的元素数量。
	// cap: 切片的容量表示在底层数组中，从切片起始位置到数组末尾的空间大小
	c, l := cap(src), len(src)
	n, changed := calChanged(c, l)
	if !changed {
		return src
	}
	// 创建一个长度为 0、容量为 n 的切片
	s := make([]T, 0, n)
	s = append([]T(nil), src...)
	return s
}

// calChanged 缩容算法
func calChanged(c, l int) (int, bool) {
	// 如果切片不大，直接返回，不用缩容
	if c <= 64 {
		return c, false
	}

	if c > 2048 && (c/l >= 2) {
		factor := 0.7
		//要缩小多少倍
		return int(float32(c) * float32(factor)), true
	}

	if c <= 2048 && (c/l >= 4) {
		return c / 2, true
	}
	return c, false
}
