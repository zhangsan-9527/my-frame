package unsafe

import "testing"

func TestPrintFieldOffest(t *testing.T) {
	testCases := []struct {
		name   string
		entity any
	}{
		{
			name:   "user",
			entity: User{},
		},
		{
			name:   "userv1",
			entity: UserV1{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			PrintFieldOffest(tc.entity)
		})
	}
}

type User struct {
	// 0
	Name string
	// 16
	Age int32
	// 24
	Alias []string
	// 48
	Address string
}

type UserV1 struct {
	// 0
	Name string
	// 16
	Age int32
	// 20
	Agev1 int32
	// 24
	Alias []string
	// 48
	Address string
}

/*
	Go unsafe - Go 对齐规则
			按照字长对齐。
				因为 Go 本身每一次访问内存都是按照字长的倍数来访问的。
					在32位字长机器上，就是按照4 个字节对产
					在64 位字长机器上，就是按照8个字节对齐

*/
