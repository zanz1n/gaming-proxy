package main

import (
	"context"
	"fmt"
	"slices"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/dns"
	"github.com/cloudflare/cloudflare-go/v6/option"
	"github.com/cloudflare/cloudflare-go/v6/zones"
)

const TAG_NAME = "github.com/zanz1n/gaming-proxy"

func RunCloudflare(ctx context.Context, host string, port uint16) error {
	cfg := GetConfig()

	c := cloudflare.NewClient(option.WithAPIToken(cfg.Cloudflare.Token))

	zone, err := c.Zones.Get(ctx, zones.ZoneGetParams{
		ZoneID: cloudflare.F(cfg.Cloudflare.ZoneID),
	})
	if err != nil {
		return err
	}

	s := state{
		Context: ctx,
		Host:    host,
		Port:    port,
		Zone:    zone,
		cfg:     &cfg.Cloudflare,
		c:       c,
		records: make([]state_record, 0),
	}

	if err = s.loadRecords(); err != nil {
		return err
	}

	if err = s.deleteTagged(); err != nil {
		return err
	}

	if err = s.recordCNAME(); err != nil {
		return err
	}

	if err = s.recordSRV(); err != nil {
		return err
	}

	return nil
}

type state_record struct {
	id      string
	name    string
	comment string
}

type state struct {
	Context context.Context
	Host    string
	Port    uint16
	Zone    *zones.Zone

	cfg     *CloudflareConfig
	c       *cloudflare.Client
	records []state_record
}

func (s *state) loadRecords() error {
	res := s.c.DNS.Records.ListAutoPaging(s.Context,
		dns.RecordListParams{
			ZoneID: cloudflare.F(s.Zone.ID),
		},
	)

	for res.Next() {
		if err := res.Err(); err != nil {
			return err
		}
		page := res.Current()

		s.records = append(s.records, state_record{
			id:      page.ID,
			name:    page.Name,
			comment: page.Comment,
		})
	}
	return nil
}

func (s *state) recordExists(name string) string {
	for _, v := range s.records {
		if v.name == name {
			return v.id
		}
	}
	return ""
}

func (s *state) deleteTagged() error {
	indexes := make([]int, 0, 2)
	for i, v := range s.records {
		if v.comment == TAG_NAME {
			fmt.Println(i)
			indexes = append(indexes, i)
		}
	}

	if len(indexes) == 0 {
		return nil
	}
	defer func() {
		s.records = s.records[:len(s.records)-len(indexes)]
		s.records = slices.Clip(s.records)
	}()

	for j, i := range indexes {
		_, err := s.c.DNS.Records.Delete(s.Context,
			s.records[i].id,
			dns.RecordDeleteParams{ZoneID: cloudflare.F(s.Zone.ID)},
		)
		if err != nil {
			return err
		}

		s.records[i] = s.records[len(s.records)-(j+1)]
	}

	return nil
}

type dataType interface {
	dns.RecordNewParamsBodyUnion
	dns.RecordUpdateParamsBodyUnion
}

func (s *state) record(name string, data dataType) (err error) {
	id := s.recordExists(name)

	if id == "" {
		_, err = s.c.DNS.Records.New(s.Context, dns.RecordNewParams{
			ZoneID: cloudflare.F(s.cfg.ZoneID),
			Body:   data,
		})
	} else {
		if !s.cfg.Overwrite {
			return fmt.Errorf(
				"DNS record `%s` already exists.\n"+
					"Delete it or set CLOUDFLARE_OVERWRITE to true",
				name,
			)
		}

		_, err = s.c.DNS.Records.Update(s.Context, id, dns.RecordUpdateParams{
			ZoneID: cloudflare.F(s.cfg.ZoneID),
			Body:   data,
		})
		return err
	}
	return
}

func (s *state) recordCNAME() error {
	name := fmt.Sprintf("%s.%s", s.cfg.Subdomain, s.Zone.Name)

	data := dns.CNAMERecordParam{
		Type:    cloudflare.F(dns.CNAMERecordTypeCNAME),
		Name:    cloudflare.F(name),
		TTL:     cloudflare.F(dns.TTL1),
		Proxied: cloudflare.F(false),
		Comment: cloudflare.F(TAG_NAME),

		Content: cloudflare.F(s.Host),
	}

	return s.record(name, data)
}

func (s *state) recordSRV() error {
	name := fmt.Sprintf("_%s._%s.%s.%s",
		s.cfg.Service,
		s.cfg.Protocol,
		s.cfg.Subdomain,
		s.Zone.Name,
	)

	data := dns.SRVRecordParam{
		Type:    cloudflare.F(dns.SRVRecordTypeSRV),
		Name:    cloudflare.F(name),
		TTL:     cloudflare.F(dns.TTL1),
		Proxied: cloudflare.F(false),
		Comment: cloudflare.F(TAG_NAME),

		Data: cloudflare.F(dns.SRVRecordDataParam{
			Port:     cloudflare.F(float64(s.Port)),
			Priority: cloudflare.F(float64(1)),
			Target: cloudflare.F(fmt.Sprintf("%s.%s",
				s.cfg.Subdomain,
				s.Zone.Name,
			)),
			Weight: cloudflare.F(float64(0)),
		}),
	}

	return s.record(name, data)
}
