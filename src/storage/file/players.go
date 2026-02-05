// Package file 文件存储角色管理
package file

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	storage "github.com/httye/yggdrasil-skins-go/src/storage/interface"
	"github.com/httye/yggdrasil-skins-go/src/yggdrasil"

	"github.com/bytedance/sonic"
)

// loadPlayers 加载角色数据
func (s *Storage) loadPlayers() error {
	playersFile := filepath.Join(s.dataDir, "players.json")

	// 如果文件不存在，创建默认数据
	if _, err := os.Stat(playersFile); os.IsNotExist(err) {
		return s.createDefaultPlayers()
	}

	data, err := os.ReadFile(playersFile)
	if err != nil {
		return err
	}

	var players []*FilePlayer
	if err := sonic.Unmarshal(data, &players); err != nil {
		return err
	}

	// 加载到缓�?	for _, player := range players {
		s.players[player.UUID] = player
		// 更新用户角色映射
		for email, user := range s.users {
			if user.UID == player.UID {
				s.userProfiles[email] = append(s.userProfiles[email], player.UUID)
				break
			}
		}
	}

	return nil
}

// savePlayers 保存角色数据
func (s *Storage) savePlayers() error {
	var players []*FilePlayer
	for _, player := range s.players {
		players = append(players, player)
	}

	data, err := sonic.MarshalIndent(players, "", "  ")
	if err != nil {
		return err
	}

	playersFile := filepath.Join(s.dataDir, "players.json")
	return os.WriteFile(playersFile, data, 0644)
}

// createDefaultPlayers 创建默认角色数据
func (s *Storage) createDefaultPlayers() error {
	// 创建测试角色
	testPlayers := []*FilePlayer{
		{
			PID:        1,
			UID:        1,
			Name:       "TestPlayer",
			UUID:       "550e8400-e29b-41d4-a716-446655440000",
			SkinTID:    0,
			CapeTID:    0,
			LastModify: time.Now().Format("2006-01-02 15:04:05"),
		},
		{
			PID:        2,
			UID:        2,
			Name:       "User2Player",
			UUID:       "550e8400-e29b-41d4-a716-446655440001",
			SkinTID:    0,
			CapeTID:    0,
			LastModify: time.Now().Format("2006-01-02 15:04:05"),
		},
		{
			PID:        3,
			UID:        3,
			Name:       "AdminPlayer",
			UUID:       "550e8400-e29b-41d4-a716-446655440002",
			SkinTID:    0,
			CapeTID:    0,
			LastModify: time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	for _, player := range testPlayers {
		s.players[player.UUID] = player
	}

	return s.savePlayers()
}

// GetProfileByUUID 根据UUID获取角色
func (s *Storage) GetProfileByUUID(uuid string) (*yggdrasil.Profile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if player, exists := s.players[uuid]; exists {
		// 获取角色的材质信�?		textures, err := s.GetPlayerTextures(uuid)
		if err != nil {
			// 如果获取材质失败，仍然返回角色信息，但properties为空
			return &yggdrasil.Profile{
				ID:         player.UUID,
				Name:       player.Name,
				Properties: []yggdrasil.ProfileProperty{},
			}, nil
		}

		// 提取皮肤和披风URL
		var skinURL, capeURL string
		var isSlim bool

		if skinInfo, exists := textures[storage.TextureTypeSkin]; exists {
			skinURL = skinInfo.URL
			if skinInfo.Metadata != nil {
				isSlim = skinInfo.Metadata.Slim
			}
		}

		if capeInfo, exists := textures[storage.TextureTypeCape]; exists {
			capeURL = capeInfo.URL
		}

		// 生成properties
		properties, err := yggdrasil.GenerateProfileProperties(player.UUID, player.Name, skinURL, capeURL, isSlim)
		if err != nil {
			// 如果生成properties失败，返回空properties
			properties = []yggdrasil.ProfileProperty{}
		}

		return &yggdrasil.Profile{
			ID:         player.UUID,
			Name:       player.Name,
			Properties: properties,
		}, nil
	}

	return nil, fmt.Errorf("profile not found")
}

// GetProfileByName 根据角色名获取角�?func (s *Storage) GetProfileByName(name string) (*yggdrasil.Profile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, player := range s.players {
		if player.Name == name {
			// 获取角色的材质信�?			textures, err := s.GetPlayerTextures(player.UUID)
			if err != nil {
				// 如果获取材质失败，仍然返回角色信息，但properties为空
				return &yggdrasil.Profile{
					ID:         player.UUID,
					Name:       player.Name,
					Properties: []yggdrasil.ProfileProperty{},
				}, nil
			}

			// 提取皮肤和披风URL
			var skinURL, capeURL string
			var isSlim bool

			if skinInfo, exists := textures[storage.TextureTypeSkin]; exists {
				skinURL = skinInfo.URL
				if skinInfo.Metadata != nil {
					isSlim = skinInfo.Metadata.Slim
				}
			}

			if capeInfo, exists := textures[storage.TextureTypeCape]; exists {
				capeURL = capeInfo.URL
			}

			// 生成properties
			properties, err := yggdrasil.GenerateProfileProperties(player.UUID, player.Name, skinURL, capeURL, isSlim)
			if err != nil {
				// 如果生成properties失败，返回空properties
				properties = []yggdrasil.ProfileProperty{}
			}

			return &yggdrasil.Profile{
				ID:         player.UUID,
				Name:       player.Name,
				Properties: properties,
			}, nil
		}
	}

	return nil, fmt.Errorf("profile not found")
}

// GetProfilesByNames 根据名称列表批量获取角色
func (s *Storage) GetProfilesByNames(names []string) ([]*yggdrasil.Profile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var profiles []*yggdrasil.Profile
	for _, name := range names {
		for _, player := range s.players {
			if player.Name == name {
				profiles = append(profiles, &yggdrasil.Profile{
					ID:         player.UUID,
					Name:       player.Name,
					Properties: []yggdrasil.ProfileProperty{}, // 初始化为空数组而不是nil
				})
				break
			}
		}
	}

	return profiles, nil
}

// GetProfilesByUserEmail 获取用户的所有角�?func (s *Storage) GetProfilesByUserEmail(userEmail string) ([]*yggdrasil.Profile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 获取用户信息
	user, exists := s.users[userEmail]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	// 获取该用户的所有角�?	var profiles []*yggdrasil.Profile
	for _, player := range s.players {
		if player.UID == user.UID {
			profiles = append(profiles, &yggdrasil.Profile{
				ID:         player.UUID,
				Name:       player.Name,
				Properties: []yggdrasil.ProfileProperty{}, // 初始化为空数组而不是nil
			})
		}
	}

	return profiles, nil
}

// CreateProfile 创建角色
func (s *Storage) CreateProfile(userEmail string, profile *yggdrasil.Profile) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查角色名是否已存�?	for _, player := range s.players {
		if player.Name == profile.Name {
			return fmt.Errorf("profile name already exists")
		}
	}

	// 检查UUID是否已存�?	if _, exists := s.players[profile.ID]; exists {
		return fmt.Errorf("profile UUID already exists")
	}

	// 获取用户信息
	user, exists := s.users[userEmail]
	if !exists {
		return fmt.Errorf("user not found")
	}

	// 创建新角�?	newPlayer := &FilePlayer{
		PID:        len(s.players) + 1, // 简单的ID生成
		UID:        user.UID,
		Name:       profile.Name,
		UUID:       profile.ID,
		SkinTID:    0,
		CapeTID:    0,
		LastModify: time.Now().Format("2006-01-02 15:04:05"),
	}

	s.players[profile.ID] = newPlayer
	s.userProfiles[userEmail] = append(s.userProfiles[userEmail], profile.ID)

	return s.savePlayers()
}

// UpdateProfile 更新角色
func (s *Storage) UpdateProfile(profile *yggdrasil.Profile) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	player, exists := s.players[profile.ID]
	if !exists {
		return fmt.Errorf("profile not found")
	}

	// 检查新名称是否与其他角色冲�?	for uuid, p := range s.players {
		if uuid != profile.ID && p.Name == profile.Name {
			return fmt.Errorf("profile name already exists")
		}
	}

	player.Name = profile.Name
	player.LastModify = time.Now().Format("2006-01-02 15:04:05")

	return s.savePlayers()
}

// DeleteProfile 删除角色
func (s *Storage) DeleteProfile(uuid string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.players[uuid]; !exists {
		return fmt.Errorf("profile not found")
	}

	// 从用户角色映射中删除
	for email, profileIDs := range s.userProfiles {
		for i, profileID := range profileIDs {
			if profileID == uuid {
				s.userProfiles[email] = append(profileIDs[:i], profileIDs[i+1:]...)
				break
			}
		}
	}

	delete(s.players, uuid)
	return s.savePlayers()
}

// ListProfiles 列出所有角色（分页�?func (s *Storage) ListProfiles(offset, limit int) ([]*yggdrasil.Profile, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var profiles []*yggdrasil.Profile
	for _, player := range s.players {
		profiles = append(profiles, &yggdrasil.Profile{
			ID:   player.UUID,
			Name: player.Name,
		})
	}

	total := len(profiles)

	// 分页处理
	if offset >= total {
		return []*yggdrasil.Profile{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return profiles[offset:end], total, nil
}
