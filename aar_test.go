package aar

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// TestAAR_ReadWrite 测试基本的读写功能
func TestAAR_ReadWrite(t *testing.T) {
	// 生成一个唯一的名称，避免测试间冲突
	name := fmt.Sprintf("test-read-write-%d", time.Now().UnixNano())
	aar, err := New(name)
	if err != nil {
		t.Fatalf("创建 AAR 实例失败: %v", err)
	}
	defer os.Remove(aar.filename) // 测试结束后清理缓存文件

	// 写入内容
	content := "hello"
	if err := aar.Write([]byte(content)); err != nil {
		t.Fatalf("写入缓存失败: %v", err)
	}

	// 立刻读取，应该成功
	// 设置一个较长的有效期
	readContent, err := aar.SetDuration(10 * time.Minute).Read()
	if err != nil {
		t.Fatalf("读取缓存失败: %v", err)
	}
	if readContent != content {
		t.Errorf("读取内容不匹配, 期望值: %s, 实际值: %s", content, readContent)
	}
}

// TestAAR_Expired_With_SetDuration 测试使用 SetDuration 设置的缓存过期
func TestAAR_Expired_With_SetDuration(t *testing.T) {
	name := fmt.Sprintf("test-expired-duration-%d", time.Now().UnixNano())
	aar, err := New(name)
	if err != nil {
		t.Fatalf("创建 AAR 实例失败: %v", err)
	}
	defer os.Remove(aar.filename)

	// 写入内容并设置一个已过期的时长（0分钟）
	content := "this should expire"
	if err := aar.Write([]byte(content)); err != nil {
		t.Fatalf("写入缓存失败: %v", err)
	}

	// 因为过期时间为0，所以应该读取失败
	_, err = aar.SetDuration(0 * time.Minute).Read()
	if err == nil {
		t.Error("期望读取失败，但实际成功了")
	}
}

// TestAAR_NotExpired_With_SetDuration 测试使用 SetDuration 设置的缓存未过期
func TestAAR_NotExpired_With_SetDuration(t *testing.T) {
	name := fmt.Sprintf("test-not-expired-duration-%d", time.Now().UnixNano())
	aar, err := New(name)
	if err != nil {
		t.Fatalf("创建 AAR 实例失败: %v", err)
	}
	defer os.Remove(aar.filename)

	// 写入内容并设置一个1分钟的有效期
	content := "this should not expire"
	if err := aar.Write([]byte(content)); err != nil {
		t.Fatalf("写入缓存失败: %v", err)
	}

	// 应该能成功读取
	readContent, err := aar.SetDuration(1 * time.Minute).Read()
	if err != nil {
		t.Fatalf("期望读取成功，但实际失败了: %v", err)
	}
	if readContent != content {
		t.Errorf("读取内容不匹配, 期望值: %s, 实际值: %s", content, readContent)
	}
}

// TestAAR_Expired_With_SetExpiredTime 测试使用 SetExpiredTime 设置的缓存过期
func TestAAR_Expired_With_SetExpiredTime(t *testing.T) {
	name := fmt.Sprintf("test-expired-time-%d", time.Now().UnixNano())
	aar, err := New(name)
	if err != nil {
		t.Fatalf("创建 AAR 实例失败: %v", err)
	}
	defer os.Remove(aar.filename)

	// 写入内容
	content := "this should expire with time"
	if err := aar.Write([]byte(content)); err != nil {
		t.Fatalf("写入缓存失败: %v", err)
	}

	// 设置一个过去的过期时间点
	expiredTime := time.Now().Add(-1 * time.Second)
	_, err = aar.SetExpiredTime(expiredTime).Read()
	if err == nil {
		t.Error("期望因过期而读取失败，但实际成功了")
	}
}

// TestAAR_NotExpired_With_SetExpiredTime 测试使用 SetExpiredTime 设置的缓存未过期
func TestAAR_NotExpired_With_SetExpiredTime(t *testing.T) {
	name := fmt.Sprintf("test-not-expired-time-%d", time.Now().UnixNano())
	aar, err := New(name)
	if err != nil {
		t.Fatalf("创建 AAR 实例失败: %v", err)
	}
	defer os.Remove(aar.filename)

	// 写入内容
	content := "this should not expire with time"
	if err := aar.Write([]byte(content)); err != nil {
		t.Fatalf("写入缓存失败: %v", err)
	}

	// 设置一个未来的过期时间点
	notExpiredTime := time.Now().Add(1 * time.Minute)
	readContent, err := aar.SetExpiredTime(notExpiredTime).Read()
	if err != nil {
		t.Fatalf("期望读取成功，但实际失败了: %v", err)
	}
	if readContent != content {
		t.Errorf("读取内容不匹配, 期望值: %s, 实际值: %s", content, readContent)
	}
}

// TestNewWithArgs 测试 New 函数支持格式化参数
func TestNewWithArgs(t *testing.T) {
	aar1, err := New("user_token_%d", 100)
	if err != nil {
		t.Fatalf("aar1 创建失败: %v", err)
	}
	defer os.Remove(aar1.filename)

	aar2, err := New("user_token_100")
	if err != nil {
		t.Fatalf("aar2 创建失败: %v", err)
	}
	defer os.Remove(aar2.filename)

	// 两个文件名应该相同
	if aar1.filename != aar2.filename {
		t.Errorf("期望文件名相同，但实际不同: %s vs %s", aar1.filename, aar2.filename)
	}
}

// TestReadNonExistent 测试读取一个不存在的缓存
func TestReadNonExistent(t *testing.T) {
	name := fmt.Sprintf("test-non-existent-%d", time.Now().UnixNano())
	aar, err := New(name)
	if err != nil {
		t.Fatalf("创建 AAR 实例失败: %v", err)
	}
	// 注意：这里没有 defer os.Remove，因为文件本身就不应该被创建

	// 直接读取一个从未写入过的缓存，应该失败
	_, err = aar.SetDuration(1 * time.Minute).Read()
	if err == nil {
		t.Error("期望读取失败，但实际成功了")
	}
}

// TestReadNonExistent 测试读取缓存内容
func TestRead(t *testing.T) {
	name := fmt.Sprintf("token-%d", time.Now().UnixNano())
	aar, err := New(name)
	if err != nil {
		t.Fatalf("创建 AAR 实例失败: %v", err)
	}
	defer os.Remove(aar.filename)

	// 直接读取一个从未写入过的缓存，应该失败
	content, err := aar.SetDuration(1 * time.Minute).Read()
	if err == nil {
		t.Error("期望读取失败，但实际成功了")
	}
	if content != "" {
		t.Errorf("期望读取到空数据，但实际为 %s", content)
	}
	err = aar.Write([]byte("abcd"))
	if err != nil {
		t.Fatalf("写入缓存失败: %v", err)
	}
	content, err = aar.SetDuration(1 * time.Minute).Read()
	if err != nil {
		t.Fatalf("读取缓存失败: %v", err)
	}
	if content != "abcd" {
		t.Errorf("读取内容不匹配, 期望值: abcd, 实际值: %s", content)
	}
}
