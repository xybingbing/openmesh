include $(TOPDIR)/rules.mk

PKG_NAME:=openmesh
PKG_VERSION:=0.1.0
PKG_RELEASE:=1

PKG_BUILD_DIR:=$(BUILD_DIR)/$(PKG_NAME)

include $(INCLUDE_DIR)/package.mk
include $(INCLUDE_DIR)/golang-package.mk

GO_PKG:=github.com/xybingbing/openmesh
GO_PKG_BUILD_PKG:=github.com/xybingbing/openmesh/cmd/openmesh

TARGET_LDFLAGS += -s -w

define Package/openmesh
  SECTION:=net
  CATEGORY:=Network
  TITLE:=OpenMesh agent
  DEPENDS:=$(GO_ARCH_DEPENDS) +wireguard-tools
endef

define Package/openmesh/description
  OpenMesh is an OpenWrt-first mesh VPN agent built around WireGuard.
endef

define Build/Prepare
	mkdir -p $(PKG_BUILD_DIR)
	$(CP) ../../* $(PKG_BUILD_DIR)/
endef

define Package/openmesh/install
	$(INSTALL_DIR) $(1)/usr/bin
	$(INSTALL_BIN) $(GO_PKG_BUILD_BIN_DIR)/openmesh $(1)/usr/bin/openmesh
	$(INSTALL_DIR) $(1)/etc/init.d
	$(INSTALL_BIN) ./files/etc/init.d/openmesh $(1)/etc/init.d/openmesh
	$(INSTALL_DIR) $(1)/etc/config
	$(INSTALL_CONF) ./files/etc/config/openmesh $(1)/etc/config/openmesh
endef

$(eval $(call GoBinPackage,openmesh))
$(eval $(call BuildPackage,openmesh))
