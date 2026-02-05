// Package blessing_skin BlessingSkin角色管理
package blessing_skin

import (
	"errors"
	"fmt"

	storage "github.com/httye/yggdrasil-skins-go/src/storage/interface"
	"github.com/httye/yggdrasil-skins-go/src/yggdrasil"

	"gorm.io/gorm"
)

// GetProfileByUUID 根据UUID获取角色（单查询优化版）
func (s *Storage) GetProfileByUUID(uuid string) (*yggdrasil.Profile, error) {
	// 一次性查询UUID映射和角色信�?	var result struct {
		PlayerName string `gorm:"column:player_name"`
		UUID       string `gorm:"column:uuid"`
	}

	err := s.db.Table("uuid u").
		Select("p.name as player_name, u.uuid").
		Joins("JOIN players p ON u.name = p.name").
		Where("u.uuid = ?", uuid).
		Take(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("profile not found")
		}
		return nil, err
	}

	// 获取角色的材质信�?	textures, err := s.GetPlayerTextures(result.UUID)
	if err != nil {
		// 如果获取材质失败，仍然返回角色信息，但properties为空
		return &yggdrasil.Profile{
			ID:         result.UUID,
			Name:       result.PlayerName,
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
	properties, err := yggdrasil.GenerateProfileProperties(result.UUID, result.PlayerName, skinURL, capeURL, isSlim)
	if err != nil {
		// 如果生成properties失败，返回空properties
		properties = []yggdrasil.ProfileProperty{}
	}

	return &yggdrasil.Profile{
		ID:         result.UUID,
		Name:       result.PlayerName,
		Properties: properties,
	}, nil
}

// GetProfileByName 根据名称获取角色（单查询优化版）
func (s *Storage) GetProfileByName(name string) (*yggdrasil.Profile, error) {
	// 一次性查询角色信息和UUID映射
	var result struct {
		PlayerName string `gorm:"column:name"`
		UUID       string `gorm:"column:uuid"`
	}

	err := s.db.Table("players p").
		Select("p.name, u.uuid").
		Joins("LEFT JOIN uuid u ON p.name = u.name").
		Where("p.name = ?", name).
		Take(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("profile not found")
		}
		return nil, err
	}

	// 如果UUID不存在，创建�?	uuid := result.UUID
	if uuid == "" {
		uuid, err = s.uuidGen.GetOrCreateUUID(name)
		if err != nil {
			return nil, err
		}
	}

	// 获取角色的材质信�?	textures, err := s.GetPlayerTextures(uuid)
	if err != nil {
		// 如果获取材质失败，仍然返回角色信息，但properties为空
		return &yggdrasil.Profile{
			ID:         uuid,
			Name:       result.PlayerName,
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
	properties, err := yggdrasil.GenerateProfileProperties(uuid, result.PlayerName, skinURL, capeURL, isSlim)
	if err != nil {
		// 如果生成properties失败，返回空properties
		properties = []yggdrasil.ProfileProperty{}
	}

	return &yggdrasil.Profile{
		ID:         uuid,
		Name:       result.PlayerName,
		Properties: properties,
	}, nil
}

// GetProfilesByNames 根据名称列表批量获取角色（优化版，自动创建UUID�?func (s *Storage) GetProfilesByNames(names []string) ([]*yggdrasil.Profile, error) {
	if len(names) == 0 {
		return []*yggdrasil.Profile{}, nil
	}

	// 1. 批量查询角色是否存在
	var players []Player
	err := s.db.Where("name IN ?", names).Find(&players).Error
	if err != nil {
		return nil, err
	}

	if len(players) == 0 {
		return []*yggdrasil.Profile{}, nil
	}

	// 2. 提取存在的角色名列表
	existingNames := make([]string, len(players))
	for i, player := range players {
		existingNames[i] = player.Name
	}

	// 3. 批量获取或创建UUID映射
	uuidMap, err := s.uuidGen.GetUUIDsByNames(existingNames)
	if err != nil {
		return nil, err
	}

	// 4. 构建结果（所有存在的角色都应该有UUID�?	var profiles []*yggdrasil.Profile
	for _, player := range players {
		if uuid, exists := uuidMap[player.Name]; exists {
			profiles = append(profiles, &yggdrasil.Profile{
				ID:         uuid,
				Name:       player.Name,
				Properties: []yggdrasil.ProfileProperty{}, // 初始化为空数组而不是nil
			})
		} else {
			// 这种情况理论上不应该发生，因为GetUUIDsByNames会自动创�?			// 但为了安全起见，我们单独处理
			uuid, err := s.uuidGen.GetOrCreateUUID(player.Name)
			if err == nil {
				profiles = append(profiles, &yggdrasil.Profile{
					ID:         uuid,
					Name:       player.Name,
					Properties: []yggdrasil.ProfileProperty{}, // 初始化为空数组而不是nil
				})
			}
		}
	}

	return profiles, nil
}

// GetProfilesByUserEmail 获取用户的所有角色（优化版）
func (s *Storage) GetProfilesByUserEmail(userEmail string) ([]*yggdrasil.Profile, error) {
	// 获取用户
	var user User
	err := s.db.Where("email = ?", userEmail).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*yggdrasil.Profile{}, nil
		}
		return nil, err
	}

	// 获取用户的所有角�?	var players []Player
	err = s.db.Where("uid = ?", user.UID).Find(&players).Error
	if err != nil {
		return nil, err
	}

	if len(players) == 0 {
		return []*yggdrasil.Profile{}, nil
	}

	// 提取角色名列�?	playerNames := make([]string, len(players))
	for i, player := range players {
		playerNames[i] = player.Name
	}

	// 批量获取或创建UUID映射
	uuidMap, err := s.uuidGen.GetUUIDsByNames(playerNames)
	if err != nil {
		return nil, err
	}

	// 构建结果（所有角色都应该有UUID�?	var profiles []*yggdrasil.Profile
	for _, player := range players {
		if uuid, exists := uuidMap[player.Name]; exists {
			profiles = append(profiles, &yggdrasil.Profile{
				ID:         uuid,
				Name:       player.Name,
				Properties: []yggdrasil.ProfileProperty{}, // 初始化为空数组而不是nil
			})
		} else {
			// 备用方案：单独创建UUID
			uuid, err := s.uuidGen.GetOrCreateUUID(player.Name)
			if err == nil {
				profiles = append(profiles, &yggdrasil.Profile{
					ID:         uuid,
					Name:       player.Name,
					Properties: []yggdrasil.ProfileProperty{}, // 初始化为空数组而不是nil
				})
			}
		}
	}

	return profiles, nil
}

// GetPlayerByName 根据名称获取BlessingSkin Player（内部使用）
func (s *Storage) GetPlayerByName(name string) (*Player, error) {
	var player Player
	err := s.db.Preload("Skin").Preload("Cape").Where("name = ?", name).First(&player).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("player not found")
		}
		return nil, err
	}
	return &player, nil
}

// GetPlayerByUUID 根据UUID获取BlessingSkin Player（内部使用）
func (s *Storage) GetPlayerByUUID(uuid string) (*Player, error) {
	playerName, err := s.uuidGen.GetNameByUUID(uuid)
	if err != nil {
		return nil, fmt.Errorf("player not found")
	}
	return s.GetPlayerByName(playerName)
}

// GetUserProfiles 根据用户UUID获取角色
func (s *Storage) GetUserProfiles(userUUID string) ([]*yggdrasil.Profile, error) {
	// 根据UUID找到用户
	var user User
	err := s.db.Where("uuid = ?", userUUID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*yggdrasil.Profile{}, nil
		}
		return nil, err
	}

	// 获取用户的所有角�?	var players []Player
	err = s.db.Where("uid = ?", user.UID).Find(&players).Error
	if err != nil {
		return nil, err
	}

	var profiles []*yggdrasil.Profile
	for _, player := range players {
		uuid, err := s.uuidGen.GetOrCreateUUID(player.Name)
		if err != nil {
			continue
		}

		profiles = append(profiles, &yggdrasil.Profile{
			ID:         uuid,
			Name:       player.Name,
			Properties: []yggdrasil.ProfileProperty{}, // 初始化为空数组而不是nil
		})
	}

	return profiles, nil
}
