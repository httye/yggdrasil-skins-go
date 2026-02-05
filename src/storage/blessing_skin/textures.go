// Package blessing_skin BlessingSkin材质管理
package blessing_skin

import (
	"fmt"
	"strings"

	storage "github.com/httye/yggdrasil-skins-go/src/storage/interface"
)

// UploadTexture BlessingSkin存储不支持材质上�?func (s *Storage) UploadTexture(textureType storage.TextureType, playerUUID string, data []byte, metadata *storage.TextureMetadata) (*storage.TextureInfo, error) {
	return nil, fmt.Errorf("texture upload is not supported in BlessingSkin storage")
}

// GetTexture 获取材质信息
func (s *Storage) GetTexture(textureType storage.TextureType, playerUUID string) (*storage.TextureInfo, error) {
	// 根据UUID获取角色
	player, err := s.GetPlayerByUUID(playerUUID)
	if err != nil {
		return nil, fmt.Errorf("player not found")
	}

	var textureID int
	switch textureType {
	case storage.TextureTypeSkin:
		textureID = player.TIDSkin
	case storage.TextureTypeCape:
		textureID = player.TIDCape
	default:
		return nil, fmt.Errorf("unsupported texture type")
	}

	if textureID <= 0 {
		return nil, fmt.Errorf("texture not found")
	}

	// 获取材质记录
	var texture Texture
	err = s.db.First(&texture, textureID).Error
	if err != nil {
		return nil, fmt.Errorf("texture not found")
	}

	return &storage.TextureInfo{
		Type: textureType,
		URL:  s.getTextureURL(texture.Hash),
		Metadata: &storage.TextureMetadata{
			Hash:       texture.Hash,
			FileSize:   int64(texture.Size),
			UploadedAt: texture.UploadAt,
		},
	}, nil
}

// DeleteTexture BlessingSkin存储不支持材质删�?func (s *Storage) DeleteTexture(textureType storage.TextureType, playerUUID string) error {
	return fmt.Errorf("texture deletion is not supported in BlessingSkin storage")
}

// GetTextureURL 计算材质URL
func (s *Storage) GetTextureURL(textureType storage.TextureType, playerUUID string) string {
	// 根据UUID获取角色
	player, err := s.GetPlayerByUUID(playerUUID)
	if err != nil {
		return ""
	}

	var textureID int
	switch textureType {
	case storage.TextureTypeSkin:
		textureID = player.TIDSkin
	case storage.TextureTypeCape:
		textureID = player.TIDCape
	default:
		return ""
	}

	if textureID <= 0 {
		return ""
	}

	// 获取材质记录
	var texture Texture
	err = s.db.First(&texture, textureID).Error
	if err != nil {
		return ""
	}

	return s.getTextureURL(texture.Hash)
}

// IsUploadSupported BlessingSkin存储不支持材质上�?func (s *Storage) IsUploadSupported() bool {
	return false
}

// getTextureURL 获取材质URL
func (s *Storage) getTextureURL(hash string) string {
	// 如果配置了texture_base_url_override，使用全局texture.base_url配置
	if s.config.TextureBaseURLOverride && s.textureConfig != nil && s.textureConfig.BaseURL != "" {
		return fmt.Sprintf("%s/%s", strings.TrimRight(s.textureConfig.BaseURL, "/"), hash)
	}

	// 默认从BlessingSkin的options表读取site_url
	siteURL := s.optionsMgr.GetOptionWithDefault("site_url", "")
	if siteURL != "" {
		return fmt.Sprintf("%s/textures/%s", siteURL, hash)
	}

	// 最后的默认�?	return fmt.Sprintf("https://your.website/textures/%s", hash)
}

// GetTextureByHash 根据哈希获取材质（内部使用）
func (s *Storage) GetTextureByHash(hash string) (*Texture, error) {
	var texture Texture
	err := s.db.Where("hash = ?", hash).First(&texture).Error
	if err != nil {
		return nil, err
	}
	return &texture, nil
}

// GetPlayerTextures 获取角色的所有材质（优化版）
func (s *Storage) GetPlayerTextures(playerUUID string) (map[storage.TextureType]*storage.TextureInfo, error) {
	// 根据UUID获取角色�?	playerName, err := s.uuidGen.GetNameByUUID(playerUUID)
	if err != nil {
		return nil, fmt.Errorf("player not found")
	}

	// 使用JOIN查询一次性获取角色和材质信息
	var result struct {
		PID      uint   `gorm:"column:pid"`
		Name     string `gorm:"column:name"`
		TIDSkin  int    `gorm:"column:tid_skin"`
		TIDCape  int    `gorm:"column:tid_cape"`
		SkinHash string `gorm:"column:skin_hash"`
		SkinSize int    `gorm:"column:skin_size"`
		SkinType string `gorm:"column:skin_type"`
		SkinTime string `gorm:"column:skin_time"`
		CapeHash string `gorm:"column:cape_hash"`
		CapeSize int    `gorm:"column:cape_size"`
		CapeTime string `gorm:"column:cape_time"`
	}

	err = s.db.Table("players p").
		Select(`p.pid, p.name, p.tid_skin, p.tid_cape,
			s.hash as skin_hash, s.size as skin_size, s.type as skin_type, s.upload_at as skin_time,
			c.hash as cape_hash, c.size as cape_size, c.upload_at as cape_time`).
		Joins("LEFT JOIN textures s ON p.tid_skin = s.tid AND p.tid_skin > 0").
		Joins("LEFT JOIN textures c ON p.tid_cape = c.tid AND p.tid_cape > 0").
		Where("p.name = ?", playerName).
		First(&result).Error

	if err != nil {
		return nil, fmt.Errorf("player not found")
	}

	textures := make(map[storage.TextureType]*storage.TextureInfo)

	// 处理皮肤材质
	if result.TIDSkin > 0 && result.SkinHash != "" {
		// 判断是否为纤细模�?		isSlim := result.SkinType == "alex"

		textures[storage.TextureTypeSkin] = &storage.TextureInfo{
			Type: storage.TextureTypeSkin,
			URL:  s.getTextureURL(result.SkinHash),
			Metadata: &storage.TextureMetadata{
				Hash:     result.SkinHash,
				FileSize: int64(result.SkinSize),
				Slim:     isSlim,
				// UploadedAt: 需要解析时间字符串，这里简化处�?			},
		}
	}

	// 处理披风材质
	if result.TIDCape > 0 && result.CapeHash != "" {
		textures[storage.TextureTypeCape] = &storage.TextureInfo{
			Type: storage.TextureTypeCape,
			URL:  s.getTextureURL(result.CapeHash),
			Metadata: &storage.TextureMetadata{
				Hash:     result.CapeHash,
				FileSize: int64(result.CapeSize),
				// UploadedAt: 需要解析时间字符串，这里简化处�?			},
		}
	}

	return textures, nil
}
