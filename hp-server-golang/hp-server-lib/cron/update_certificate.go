package cron

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"hp-server-lib/db"
	"hp-server-lib/entity"
	"hp-server-lib/net/acme"
	"log"
	"time"

	"github.com/gogf/gf/v2/os/gcron"
)

// 初始化证书更新任务，每日执行
func initUpdateCertificate(ctx context.Context) {
	gcron.AddSingleton(ctx, "@daily", updateCertificate, "update_certificate")
}

// 更新证书
func updateCertificate(ctx context.Context) {
	// 获取所有域名
	var results []*entity.UserDomainEntity
	tx := db.DB.Model(&entity.UserDomainEntity{}).Find(&results)
	if tx.Error != nil {
		log.Printf("数据库查询失败: %v", tx.Error)
		return
	}
	for _, data := range results {
		if len(data.CertificateKey) > 0 && len(data.CertificateContent) > 0 {
			// 解码 PEM 数据
			block, _ := pem.Decode([]byte(data.CertificateContent))
			if block == nil {
				log.Printf("PEM 数据无效")
				continue
			}

			// 解析证书
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				log.Printf("证书解析失败: %v", err)
				continue
			}

			// 检查证书包含的域名
			if len(cert.DNSNames) == 0 {
				log.Printf("证书未包含任何域名")
				continue
			}

			// 判断证书剩余时间
			if time.Until(cert.NotAfter) < time.Hour*24*30 {
				log.Printf("证书 %v 即将过期: %s", cert.DNSNames, cert.NotAfter.String())
				// 重新申请证书
				cert, err := acme.ConfigAcme.GenCert(cert.DNSNames[0])
				if err != nil {
					log.Printf("证书更新失败: %v", err)
					continue
				}
				log.Printf("证书获取成功")
				tx := db.DB.Model(&entity.UserDomainEntity{}).
					Where("id = ?", data.Id).
					Updates(map[string]interface{}{
						"certificate_key":     string(cert.PrivateKey),
						"certificate_content": string(cert.Certificate),
					})
				if tx.Error != nil {
					log.Printf("数据库更新失败: %v", tx.Error)
				} else {
					log.Printf("证书更新成功")
				}
			}
		}
	}
}
