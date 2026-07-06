package nginx

import (
	"context"

	"github.com/Ajayvtl/devserver/internal/logger"
)

type Module struct{}

func (m *Module) Name() string { return "nginx" }

func (m *Module) Description() string { return "Web server and reverse proxy" }

func (m *Module) Check(ctx context.Context) error     { return m.log(ctx, "check") }
func (m *Module) Install(ctx context.Context) error   { return m.log(ctx, "install") }
func (m *Module) Configure(ctx context.Context) error { return m.log(ctx, "configure") }
func (m *Module) Validate(ctx context.Context) error  { return m.log(ctx, "validate") }
func (m *Module) Upgrade(ctx context.Context) error   { return m.log(ctx, "upgrade") }
func (m *Module) Uninstall(ctx context.Context) error { return m.log(ctx, "uninstall") }
func (m *Module) Rollback(ctx context.Context) error  { return m.log(ctx, "rollback") }

func (m *Module) log(ctx context.Context, action string) error {
	logger.FromContext(ctx).Info().Str("module", m.Name()).Str("action", action).Msg("module step")
	return nil
}
