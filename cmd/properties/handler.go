package main

import (
	"context"
	"log"
	"time"
)

type handler struct {
	user    UserStore
	setting SettingStore
}

// Get returs list of settings names and values, for given user and period of time.
func (h *handler) GetSettings(ctx context.Context, userID int, period time.Time) ([]Setting, error) {
	log.Printf("get-settings for %d at '%s'", userID, period)

	us, err := h.user.Get(ctx, userID, period)
	if err != nil {
		return nil, err
	}

	return h.setting.Get(ctx, period, us.Bundles)
}

// ListSettings returns list of settings names.
func (h *handler) ListSettings(ctx context.Context) ([]string, error) {
	return h.setting.SettingsList(ctx)
}

// ListBundles returns list of bundles names.
func (h *handler) ListBundles(ctx context.Context) ([]Bundle, error) {
	return h.setting.BundlesList(ctx)
}

// ListTags returns list of existing tags.
func (h *handler) ListTags(ctx context.Context) ([]string, error) {
	return h.setting.TagsList(ctx)
}

// SetTag sets new tag for user.
func (h *handler) SetTag(ctx context.Context, userID int, tag string, expire *time.Time) error {
	us, err := h.user.Get(ctx, userID, time.Now())
	if err != nil {
		return err
	}

	curb, err := h.setting.BundlesByID(ctx, us.Bundles)
	if err != nil {
		return err
	}

	newb, err := h.setting.BundlesByTag(ctx, tag)
	if err != nil {
		return err
	}

	us.Expire = expire
	us.Bundles = MergeBundles(curb, newb)

	return h.user.Set(ctx, userID, us)
}

// SetBundles sets one or more bundles for user.
func (h *handler) SetBundles(ctx context.Context, userID int, bundles []string, expire *time.Time) error {
	us, err := h.user.Get(ctx, userID, time.Now())
	if err != nil {
		return err
	}

	curb, err := h.setting.BundlesByID(ctx, us.Bundles)
	if err != nil {
		return err
	}

	newb, err := h.setting.BundlesByName(ctx, bundles)
	if err != nil {
		return err
	}

	us.Expire = expire
	us.Bundles = MergeBundles(curb, newb)

	return h.user.Set(ctx, userID, us)
}

// UnSetTag un-sets tag for user.
func (h *handler) UnSetTag(ctx context.Context, userID int, tag string) error {
	us, err := h.user.Get(ctx, userID, time.Now())
	if err != nil {
		return err
	}

	curb, err := h.setting.BundlesByID(ctx, us.Bundles)
	if err != nil {
		return err
	}

	cutb, err := h.setting.BundlesByTag(ctx, tag)
	if err != nil {
		return err
	}

	us.Bundles = DropBundles(curb, cutb)

	return h.user.Set(ctx, userID, us)
}

// UnSetBundles un-sets bundles for user.
func (h *handler) UnSetBundles(ctx context.Context, userID int, bundles []string) error {
	us, err := h.user.Get(ctx, userID, time.Now())
	if err != nil {
		return err
	}

	curb, err := h.setting.BundlesByID(ctx, us.Bundles)
	if err != nil {
		return err
	}

	cutb, err := h.setting.BundlesByName(ctx, bundles)
	if err != nil {
		return err
	}

	us.Bundles = DropBundles(curb, cutb)

	return h.user.Set(ctx, userID, us)
}
