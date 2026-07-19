---
title: "XCP-NG states.db Recovery"
date: 2026-06-20
created: 2026-05-01 11:53
updated: 2026-06-20
summary: "My states.db for an XCP-NG host got corrupted due to unexpected power loss. Here is how I fixed it"
author: "Ethan Ashley"
tags: ["xcp-ng", "projects", "virtualization", "troubleshooting", "homelab"]
draft: false
---

# Symptoms
No network interfaces exist
![[attachments/97f61e2e.png]]
Web UI does not work
	Firefox gave me the error `PR_CONNECT_RESET_ERROR`

SSH to XCP-NG works
	`xe` commands do not work, not even `xe help`
	The error i got was something about connection refused to xapi

Errors
> [!CODE]- `# tail -n 500 /var/log/xensource.log`
> ```
> Apr 30 20:26:40 Hostname xapi: [error||0 |Xapi.watchdog|backtrace] 13/14 xapi Called from file ocaml/xapi/server_helpers.ml, line 97
> Apr 30 20:26:40 Hostname xapi: [error||0 |Xapi.watchdog|backtrace] 14/14 xapi Called from file ocaml/libs/log/debug.ml, line 258
> Apr 30 20:26:40 Hostname xapi: [error||0 |Xapi.watchdog|backtrace]
> Apr 30 20:26:40 Hostname xapi: [debug||0 |Xapi.watchdog|xapi] xapi top-level caught exception: INTERNAL_ERROR: [ Xmlm.Error(1:1, "expected root element") ]
> Apr 30 20:26:40 Hostname xapi: [error||0 ||backtrace] Xapi.watchdog failed with exception Xmlm.Error(1:1, "expected root element")
> Apr 30 20:26:40 Hostname xapi: [error||0 ||backtrace] Raised Xmlm.Error(1:1, "expected root element")
> Apr 30 20:26:40 Hostname xapi: [error||0 ||backtrace] 1/11 xapi Raised at file ocaml/libs/log/debug.ml, line 279
> Apr 30 20:26:40 Hostname xapi: [error||0 ||backtrace] 2/11 xapi Called from file ocaml/libs/xapi-stdext/lib/xapi-stdext-pervasives/pervasiveext.ml, line 24
> Apr 30 20:26:40 Hostname xapi: [error||0 ||backtrace] 3/11 xapi Called from file ocaml/libs/xapi-stdext/lib/xapi-stdext-pervasives/pervasiveext.ml, line 39
> Apr 30 20:26:40 Hostname xapi: [error||0 ||backtrace] 4/11 xapi Called from file ocaml/libs/xapi-stdext/lib/xapi-stdext-pervasives/pervasiveext.ml, line 24
> Apr 30 20:26:40 Hostname xapi: [error||0 ||backtrace] 5/11 xapi Called from file ocaml/libs/xapi-stdext/lib/xapi-stdext-pervasives/pervasiveext.ml, line 39
> Apr 30 20:26:40 Hostname xapi: [error||0 ||backtrace] 6/11 xapi Called from file ocaml/xapi/xapi.ml, line 1073
> Apr 30 20:26:40 Hostname xapi: [error||0 ||backtrace] 7/11 xapi Called from file ocaml/xapi/xapi.ml, line 1533
> Apr 30 20:26:40 Hostname xapi: [error||0 ||backtrace] 8/11 xapi Called from file ocaml/xapi/xapi.ml, line 1541
> Apr 30 20:26:40 Hostname xapi: [error||0 ||backtrace] 9/11 xapi Called from file ocaml/xapi/xapi.ml, line 1547
> Apr 30 20:26:40 Hostname xapi: [error||0 ||backtrace] 10/11 xapi Called from file ocaml/xapi/xapi.ml, line 1552
> Apr 30 20:26:40 Hostname xapi: [error||0 ||backtrace] 11/11 xapi Called from file ocaml/libs/log/debug.ml, line 258
> Apr 30 20:26:40 Hostname xapi: [error||0 ||backtrace]
> Apr 30 20:26:41 Hostname xapi: [ warn||0 ||xapi] Duplicate configuration keys in Xcp_service.configure: disable-logging-for in [ use-switch; switch-path; search-path; pidfile; log; disable-logging-for; loglevel; inventory; config; config-dir; timeslice; master_connection_reset_timeout; master_connection_retry_timeout; master_connection_default_timeout; qemu_dm_ready_timeout; hotplug_timeout; pif_reconfigure_ip_timeout; pool_db_sync_interval; pool_data_sync_interval; domain_shutdown_total_timeout; emergency_reboot_delay_base; emergency_reboot_delay_extra; ha_xapi_healthcheck_interval; ha_xapi_healthcheck_timeout; ha_xapi_restart_attempts; ha_xapi_restart_timeout; logrotate_check_interval; rrd_backup_interval; session_revalidation_interval; update_all_subjects_interval; wait_memory_target_timeout; snapshot_with_quiesce_timeout; host_heartbeat_interval; host_assumed_dead_interval; fuse_time; db_restore_fuse_time; inactive_session_timeout; pending_task_timeout; completed_task_timeout; minimum_time_between_bounces; minimum_time_between_reboot_with_no_added_delay; ha_monitor_interval; ha_monitor_plan_interval; ha_monitor_startup_timeout; ha_default_timeout_base; guest_liveness_timeout; permanent_master_failure_retry_interval; redo_log_max_block_time_empty; redo_log_max_block_time_read; redo_log_max_block_time_writedelta; redo_log_max_block_time_writedb; redo_log_max_startup_time; redo_log_connect_delay; default-vbd3-polling-duration; default-vbd3-polling-idle-threshold; vm_call_plugin_interval; xapi_clusterd_port; max_active_sr_scans; winbind_debug_level; winbind_cache_time; winbind_machine_pwd_timeout; winbind_update_closest_kdc_interval; header_read_timeout_tcp; header_total_timeout_tcp; max_header_length_tcp; coordinator_max_stunnel_cache; member_max_stunnel_cache; conn_limit_tcp; conn_limit_unix; conn_limit_clientcert; stunnel_cache_max_age; stunnel_cache_max_idle; export_interval; max_spans; max_traces; max_observer_file_size; test-open; local_yum_repo_port; ha_best_effort_max_retries; winbind_ldap_query_subject_timeout; sm-plugins; hotfix-fingerprint; logconfig; writereadyfile; writeinitcomplete; nowatchdog; log-getter; onsystemboot; relax-xsm-sr-check; include-console-username-in-error; disable-logging-for; disable-dbsync-for; xenopsd-queues; xenopsd-default; nvidia-whitelist; igd-passthru-vendor-whitelist; gvt-g-whitelist; gvt-g-supported; mxgpu-whitelist; pass-through-pif-carrier; cluster-stack-default; gpumon_stop_timeout; reboot_required_hfxs; xen_livepatch_list; kpatch_list; modprobe_path; db_idempotent_map; use-event-next; nvidia_multi_vgpu_enabled_driver_versions; nvidia_t4_sriov; create-tools-sr; allow-host-sched-gran-modification; use-xmlrpc; extauth_ad_backend; winbind_kerberos_encryption_type; winbind_set_machine_account_kerberos_encryption_type; winbind_allow_kerberos_auth_fallback; winbind_scan_trusted_domains; winbind_keep_configuration; hsts_max_age; website-https-only; migration-https-only; repository-domain-name-allowlist; repository-url-blocklist; repository-gpgcheck; repository-gpgkey-name; failed-login-alert-freq; cert-expiration-days; messages-limit; evacuation-batch-size; ignore-vtpm-unimplemented; allow-custom-uefi-certs; server-cert-group-id; export-interval; export-chunk-size; max-spans; max-traces; prefer-nbd-attach; observer-max-file-size; nvidia-gpumon-detach; compress-tracing-files; observer-endpoint-http-enabled; observer-endpoint-https-enabled; observer-experimental-components; disable-webserver; tgroups-enabled; event_from_delay; event_from_task_delay; event_next_delay; drivertool; reuse-pool-sessions; validate-reusable-pool-session; ssh-auto-mode; secure-boot-efi-path; vm-sysprep-enabled; vm-sysprep-wait; vhd-legacy-blocks-format; proxy_poll_period_timeout; max-span-depth; https-only-default; firewall-backend; dynamic-control-firewalld-service; ntp-service; ntp-config-path; ntp-dhcp-script-path; ntp-dhcp-dir; ntp-client-path; timedatectl; legacy-factory-ntp-servers; factory-ntp-servers; post-install-scripts-dir; gpg-homedir; xen-cmdline; cluster-stack-root; web-dir; sm-dir; udhcpd-skel; db-config-file; pool_config_file; udevadm; dracut; depmod; rpm-cmd; modifyrepo-cmd; createrepo-cmd; set-iscsi-initiator; openssl_path; gencert; alert-certificate-check; systemctl; list_domains; xen-cmdline-script; static-vdis; xsh; xe-toolstack-restart; xe; host-restore; host-backup; upload-wrapper; update-mh-info; logs-download; xe-syslog-reconfigure; set-hostname; host-bugreport-upload; fence; qcow_stream_tool; qcow_to_stdout; vhd-tool; sparse_dd; redo-log-block-device-io; pbis-force-domain-leave-script; busybox; startup-script-hook; rolling-upgrade-script-hook; xapi-message-script; non-managed-pifs; domain_join_cli_cmd; update-issue; killall; nbd-firewall-config; firewall-port-config; firewall-cmd; firewall-cmd-wrapper; nbd_client_manager; varstore-rm; varstore-sb-state; varstore-ls; varstore_dir; nvidia-sriov-manage; gen_pool_secret_script; samba administration tool; Samba TDB (Trivial Database) management tool; winbind query tool; SQLite database  management tool; yum-cmd; dnf-cmd; reposync-cmd; yum-config-manager-cmd; c_rehash; fcoe-driver; pvsproxy_close_cache_vdi; genisoimage; pool_secret_path; udhcpd-conf; remote-db-conf-file; logconfig; cpu-info-file; server-cert-path; server-cert-internal-path; stunnel-bundle-path; pool-bundle-path; iscsi_initiatorname; master-scripts-dir; packs-dir; xapi-hooks-root; xapi-plugins-root; xapi-extensions-root; static-vdis-root; tools-sr-dir; trusted-pool-certs-dir; trusted-certs-dir; trace-log-dir; pool-recommendations-dir ]
> Apr 30 20:26:41 Hostname xapi: [ warn||0 ||xapi] Duplicate configuration keys in Xcp_service.configure: logconfig in [ use-switch; switch-path; search-path; pidfile; log; disable-logging-for; loglevel; inventory; config; config-dir; timeslice; master_connection_reset_timeout; master_connection_retry_timeout; master_connection_default_timeout; qemu_dm_ready_timeout; hotplug_timeout; pif_reconfigure_ip_timeout; pool_db_sync_interval; pool_data_sync_interval; domain_shutdown_total_timeout; emergency_reboot_delay_base; emergency_reboot_delay_extra; ha_xapi_healthcheck_interval; ha_xapi_healthcheck_timeout; ha_xapi_restart_attempts; ha_xapi_restart_timeout; logrotate_check_interval; rrd_backup_interval; session_revalidation_interval; update_all_subjects_interval; wait_memory_target_timeout; snapshot_with_quiesce_timeout; host_heartbeat_interval; host_assumed_dead_interval; fuse_time; db_restore_fuse_time; inactive_session_timeout; pending_task_timeout; completed_task_timeout; minimum_time_between_bounces; minimum_time_between_reboot_with_no_added_delay; ha_monitor_interval; ha_monitor_plan_interval; ha_monitor_startup_timeout; ha_default_timeout_base; guest_liveness_timeout; permanent_master_failure_retry_interval; redo_log_max_block_time_empty; redo_log_max_block_time_read; redo_log_max_block_time_writedelta; redo_log_max_block_time_writedb; redo_log_max_startup_time; redo_log_connect_delay; default-vbd3-polling-duration; default-vbd3-polling-idle-threshold; vm_call_plugin_interval; xapi_clusterd_port; max_active_sr_scans; winbind_debug_level; winbind_cache_time; winbind_machine_pwd_timeout; winbind_update_closest_kdc_interval; header_read_timeout_tcp; header_total_timeout_tcp; max_header_length_tcp; coordinator_max_stunnel_cache; member_max_stunnel_cache; conn_limit_tcp; conn_limit_unix; conn_limit_clientcert; stunnel_cache_max_age; stunnel_cache_max_idle; export_interval; max_spans; max_traces; max_observer_file_size; test-open; local_yum_repo_port; ha_best_effort_max_retries; winbind_ldap_query_subject_timeout; sm-plugins; hotfix-fingerprint; logconfig; writereadyfile; writeinitcomplete; nowatchdog; log-getter; onsystemboot; relax-xsm-sr-check; include-console-username-in-error; disable-logging-for; disable-dbsync-for; xenopsd-queues; xenopsd-default; nvidia-whitelist; igd-passthru-vendor-whitelist; gvt-g-whitelist; gvt-g-supported; mxgpu-whitelist; pass-through-pif-carrier; cluster-stack-default; gpumon_stop_timeout; reboot_required_hfxs; xen_livepatch_list; kpatch_list; modprobe_path; db_idempotent_map; use-event-next; nvidia_multi_vgpu_enabled_driver_versions; nvidia_t4_sriov; create-tools-sr; allow-host-sched-gran-modification; use-xmlrpc; extauth_ad_backend; winbind_kerberos_encryption_type; winbind_set_machine_account_kerberos_encryption_type; winbind_allow_kerberos_auth_fallback; winbind_scan_trusted_domains; winbind_keep_configuration; hsts_max_age; website-https-only; migration-https-only; repository-domain-name-allowlist; repository-url-blocklist; repository-gpgcheck; repository-gpgkey-name; failed-login-alert-freq; cert-expiration-days; messages-limit; evacuation-batch-size; ignore-vtpm-unimplemented; allow-custom-uefi-certs; server-cert-group-id; export-interval; export-chunk-size; max-spans; max-traces; prefer-nbd-attach; observer-max-file-size; nvidia-gpumon-detach; compress-tracing-files; observer-endpoint-http-enabled; observer-endpoint-https-enabled; observer-experimental-components; disable-webserver; tgroups-enabled; event_from_delay; event_from_task_delay; event_next_delay; drivertool; reuse-pool-sessions; validate-reusable-pool-session; ssh-auto-mode; secure-boot-efi-path; vm-sysprep-enabled; vm-sysprep-wait; vhd-legacy-blocks-format; proxy_poll_period_timeout; max-span-depth; https-only-default; firewall-backend; dynamic-control-firewalld-service; ntp-service; ntp-config-path; ntp-dhcp-script-path; ntp-dhcp-dir; ntp-client-path; timedatectl; legacy-factory-ntp-servers; factory-ntp-servers; post-install-scripts-dir; gpg-homedir; xen-cmdline; cluster-stack-root; web-dir; sm-dir; udhcpd-skel; db-config-file; pool_config_file; udevadm; dracut; depmod; rpm-cmd; modifyrepo-cmd; createrepo-cmd; set-iscsi-initiator; openssl_path; gencert; alert-certificate-check; systemctl; list_domains; xen-cmdline-script; static-vdis; xsh; xe-toolstack-restart; xe; host-restore; host-backup; upload-wrapper; update-mh-info; logs-download; xe-syslog-reconfigure; set-hostname; host-bugreport-upload; fence; qcow_stream_tool; qcow_to_stdout; vhd-tool; sparse_dd; redo-log-block-device-io; pbis-force-domain-leave-script; busybox; startup-script-hook; rolling-upgrade-script-hook; xapi-message-script; non-managed-pifs; domain_join_cli_cmd; update-issue; killall; nbd-firewall-config; firewall-port-config; firewall-cmd; firewall-cmd-wrapper; nbd_client_manager; varstore-rm; varstore-sb-state; varstore-ls; varstore_dir; nvidia-sriov-manage; gen_pool_secret_script; samba administration tool; Samba TDB (Trivial Database) management tool; winbind query tool; SQLite database  management tool; yum-cmd; dnf-cmd; reposync-cmd; yum-config-manager-cmd; c_rehash; fcoe-driver; pvsproxy_close_cache_vdi; genisoimage; pool_secret_path; udhcpd-conf; remote-db-conf-file; logconfig; cpu-info-file; server-cert-path; server-cert-internal-path; stunnel-bundle-path; pool-bundle-path; iscsi_initiatorname; master-scripts-dir; packs-dir; xapi-hooks-root; xapi-plugins-root; xapi-extensions-root; static-vdis-root; tools-sr-dir; trusted-pool-certs-dir; trusted-certs-dir; trace-log-dir; pool-recommendations-dir ]
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] Parsing [http db_write redo_log api_readonly]
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi_globs] Whitelisting PCI vendor 8086 for passthrough
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] Parsing [tracing tracing_export]
> Apr 30 20:26:41 Hostname xapi: [ info||0 ||xapi] Found 'alert-certificate-check' at '/opt/xensource/libexec/alert-certificate-check'
> Apr 30 20:26:41 Hostname xapi: [ info||0 ||xapi] Found 'static-vdis' at '/opt/xensource/bin/static-vdis'
> Apr 30 20:26:41 Hostname xapi: [ info||0 ||xapi] Found 'xsh' at '/opt/xensource/bin/xsh'
> Apr 30 20:26:41 Hostname xapi: [ info||0 ||xapi] Found 'xe-toolstack-restart' at '/opt/xensource/bin/xe-toolstack-restart'
> Apr 30 20:26:41 Hostname xapi: [ info||0 ||xapi] Found 'xe' at '/usr/bin/xe'
> Apr 30 20:26:41 Hostname xapi: [ info||0 ||xapi] Found 'host-restore' at '/opt/xensource/libexec/host-restore'
> Apr 30 20:26:41 Hostname xapi: [ info||0 ||xapi] Found 'host-backup' at '/opt/xensource/libexec/host-backup'
> Apr 30 20:26:41 Hostname xapi: [ info||0 ||xapi] Found 'upload-wrapper' at '/opt/xensource/libexec/upload-wrapper'
> Apr 30 20:26:41 Hostname xapi: [ info||0 ||xapi] Found 'update-mh-info' at '/opt/xensource/libexec/update-mh-info'
> Apr 30 20:26:41 Hostname xapi: [ info||0 ||xapi] Found 'logs-download' at '/opt/xensource/libexec/logs-download'
> Apr 30 20:26:41 Hostname xapi: [ info||0 ||xapi] Found 'xe-syslog-reconfigure' at '/opt/xensource/libexec/xe-syslog-reconfigure'
> Apr 30 20:26:41 Hostname xapi: [ info||0 ||xapi] Found 'set-hostname' at '/opt/xensource/libexec/set-hostname'
> Apr 30 20:26:41 Hostname xapi: [ info||0 ||xapi] Found 'host-bugreport-upload' at '/opt/xensource/libexec/host-bugreport-upload'
> Apr 30 20:26:41 Hostname xapi: [ info||0 ||xapi] Found 'fence' at '/opt/xensource/libexec/fence'
> Apr 30 20:26:41 Hostname xapi: [ info||0 ||xapi] Found 'qcow-stream-tool' at '/usr/bin/qcow-stream-tool'
> Apr 30 20:26:41 Hostname xapi: [ info||0 ||xapi] Found 'vhd-tool' at '/usr/bin/vhd-tool'
> Apr 30 20:26:41 Hostname xapi: [ info||0 ||xapi] Found 'block_device_io' at '/opt/xensource/libexec/block_device_io'
> Apr 30 20:26:41 Hostname xapi: [ info||0 ||xapi] Found 'pbis-force-domain-leave' at '/opt/xensource/libexec/pbis-force-domain-leave'
> Apr 30 20:26:41 Hostname xapi: [ info||0 ||xapi] Found 'busybox' at '/usr/bin/busybox'
> Apr 30 20:26:41 Hostname xapi: [ warn||0 ||xapi] Failed to find xapi-startup-script on $PATH ( = /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin) or search_path option ( = /opt/xensource/libexec:/opt/xensource/bin)
> Apr 30 20:26:41 Hostname xapi: [ info||0 ||xapi] Found 'xapi-rolling-upgrade' at '/opt/xensource/libexec/xapi-rolling-upgrade'
> Apr 30 20:26:41 Hostname xapi: [ info||0 ||xapi] Found 'mail-alarm' at '/opt/xensource/libexec/mail-alarm'
> Apr 30 20:26:41 Hostname xapi: [ info||0 ||xapi] Found 'update-issue' at '/usr/sbin/update-issue'
> Apr 30 20:26:41 Hostname xapi: [ info||0 ||xapi] Found 'killall' at '/usr/bin/killall'
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] use-switch = true (true if the message switch is to be enabled)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] switch-path = /var/run/message-switch/sock (Unix domain socket path on localhost where the message switch is listening)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] search-path = /opt/xensource/libexec:/opt/xensource/bin (Search path for resources)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] pidfile = /var/run/xapi.pid (Filename to write process PID)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] log = syslog:daemon (Where to write log messages)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] disable-logging-for = tracing_export tracing redo_log mscgen http db_write api_readonly (A space-separated list of debug modules to suppress logging from)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] loglevel = debug (Log level)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] inventory = /etc/xensource-inventory (Location of the inventory file)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] config = /etc/xapi.conf (Location of configuration file)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] config-dir = /etc/xapi.conf.d (Location of directory containing configuration file fragments)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] timeslice = 0.050 (timeslice in seconds)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] master_connection_reset_timeout = 120. (Set the value of 'master_connection_reset_timeout')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] master_connection_retry_timeout = -1. (Set the value of 'master_connection_retry_timeout')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] master_connection_default_timeout = 10. (Set the value of 'master_connection_default_timeout')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] qemu_dm_ready_timeout = 300. (Set the value of 'qemu_dm_ready_timeout')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] hotplug_timeout = 300. (Set the value of 'hotplug_timeout')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] pif_reconfigure_ip_timeout = 300. (Set the value of 'pif_reconfigure_ip_timeout')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] pool_db_sync_interval = 300. (Set the value of 'pool_db_sync_interval')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] pool_data_sync_interval = 86400. (Set the value of 'pool_data_sync_interval')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] domain_shutdown_total_timeout = 1200. (Set the value of 'domain_shutdown_total_timeout')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] emergency_reboot_delay_base = 60. (Set the value of 'emergency_reboot_delay_base')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] emergency_reboot_delay_extra = 120. (Set the value of 'emergency_reboot_delay_extra')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] ha_xapi_healthcheck_interval = 60 (Set the value of 'ha_xapi_healthcheck_interval')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] ha_xapi_healthcheck_timeout = 120 (Set the value of 'ha_xapi_healthcheck_timeout')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] ha_xapi_restart_attempts = 1 (Set the value of 'ha_xapi_restart_attempts')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] ha_xapi_restart_timeout = 300 (Set the value of 'ha_xapi_restart_timeout')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] logrotate_check_interval = 300. (Set the value of 'logrotate_check_interval')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] rrd_backup_interval = 86400. (Set the value of 'rrd_backup_interval')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] session_revalidation_interval = 300. (Set the value of 'session_revalidation_interval')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] update_all_subjects_interval = 900. (Set the value of 'update_all_subjects_interval')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] wait_memory_target_timeout = 256. (Set the value of 'wait_memory_target_timeout')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] snapshot_with_quiesce_timeout = 600. (Set the value of 'snapshot_with_quiesce_timeout')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] host_heartbeat_interval = 30. (Set the value of 'host_heartbeat_interval')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] host_assumed_dead_interval = 600000000000ns (10min) (Set the value of 'host_assumed_dead_interval')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] fuse_time = 10. (Set the value of 'fuse_time')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] db_restore_fuse_time = 30. (Set the value of 'db_restore_fuse_time')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] inactive_session_timeout = 86400. (Set the value of 'inactive_session_timeout')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] pending_task_timeout = 86400. (Set the value of 'pending_task_timeout')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] completed_task_timeout = 3900. (Set the value of 'completed_task_timeout')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] minimum_time_between_bounces = 120. (Set the value of 'minimum_time_between_bounces')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] minimum_time_between_reboot_with_no_added_delay = 60. (Set the value of 'minimum_time_between_reboot_with_no_added_delay')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] ha_monitor_interval = 20. (Set the value of 'ha_monitor_interval')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] ha_monitor_plan_interval = 1800. (Set the value of 'ha_monitor_plan_interval')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] ha_monitor_startup_timeout = 1800. (Set the value of 'ha_monitor_startup_timeout')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] ha_default_timeout_base = 60. (Set the value of 'ha_default_timeout_base')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] guest_liveness_timeout = 300. (Set the value of 'guest_liveness_timeout')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] permanent_master_failure_retry_interval = 60. (Set the value of 'permanent_master_failure_retry_interval')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] redo_log_max_block_time_empty = 2. (Set the value of 'redo_log_max_block_time_empty')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] redo_log_max_block_time_read = 30. (Set the value of 'redo_log_max_block_time_read')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] redo_log_max_block_time_writedelta = 2. (Set the value of 'redo_log_max_block_time_writedelta')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] redo_log_max_block_time_writedb = 30. (Set the value of 'redo_log_max_block_time_writedb')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] redo_log_max_startup_time = 5. (Set the value of 'redo_log_max_startup_time')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] redo_log_connect_delay = 0.1 (Set the value of 'redo_log_connect_delay')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] default-vbd3-polling-duration = 8000 (Set the value of 'default-vbd3-polling-duration')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] default-vbd3-polling-idle-threshold = 50 (Set the value of 'default-vbd3-polling-idle-threshold')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] vm_call_plugin_interval = 10. (Set the value of 'vm_call_plugin_interval')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] xapi_clusterd_port = 8896 (Set the value of 'xapi_clusterd_port')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] max_active_sr_scans = 32 (Set the value of 'max_active_sr_scans')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] winbind_debug_level = 2 (Set the value of 'winbind_debug_level')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] winbind_cache_time = 60 (Set the value of 'winbind_cache_time')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] winbind_machine_pwd_timeout = 1209600. (Set the value of 'winbind_machine_pwd_timeout')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] winbind_update_closest_kdc_interval = 79200. (Set the value of 'winbind_update_closest_kdc_interval')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] header_read_timeout_tcp = 10. (Set the value of 'header_read_timeout_tcp')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] header_total_timeout_tcp = 60. (Set the value of 'header_total_timeout_tcp')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] max_header_length_tcp = 1024 (Set the value of 'max_header_length_tcp')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] coordinator_max_stunnel_cache = 70 (Set the value of 'coordinator_max_stunnel_cache')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] member_max_stunnel_cache = 70 (Set the value of 'member_max_stunnel_cache')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] conn_limit_tcp = 800 (Set the value of 'conn_limit_tcp')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] conn_limit_unix = 1024 (Set the value of 'conn_limit_unix')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] conn_limit_clientcert = 800 (Set the value of 'conn_limit_clientcert')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] stunnel_cache_max_age = 10800. (Set the value of 'stunnel_cache_max_age')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] stunnel_cache_max_idle = 300. (Set the value of 'stunnel_cache_max_idle')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] export_interval = 30. (Set the value of 'export_interval')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] max_spans = 10000 (Set the value of 'max_spans')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] max_traces = 10000 (Set the value of 'max_traces')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] max_observer_file_size = 1048576 (Set the value of 'max_observer_file_size')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] test-open = 0 (Set the value of 'test-open')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] local_yum_repo_port = 8000 (Set the value of 'local_yum_repo_port')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] ha_best_effort_max_retries = 2 (Set the value of 'ha_best_effort_max_retries')
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] winbind_ldap_query_subject_timeout = 20000000000ns (20s) (Timeout to perform ldap query for subject information)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] sm-plugins = [ largeblock; zfs; xfs; glusterfs; cephfs; smb; shm; lvmofcoe; lvmohba; lvm; iso; udev; rawhba; hba; file; dummy; lvmoiscsi; iscsi; nfs; ext ] (space-separated list of storage plugins to allow.)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] hotfix-fingerprint = NERDNTUzMDMwRUMwNDFFNDI4N0M4OEVCRUFEMzlGOTJEOEE5REUyNg== (Fingerprint of the key used for signed hotfixes)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] logconfig = /etc/xensource/log.conf (Log config file to use)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] writereadyfile = /var/run/xapi_startup.cookie (touch specified file when xapi is ready to accept requests)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] writeinitcomplete = /var/run/xapi_init_complete.cookie (touch specified file when xapi init process is complete)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] nowatchdog = true (turn watchdog off, avoiding initial fork)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] log-getter = false (Enable/Disable logging for getters)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] onsystemboot = false (indicates that this server start is the first since the host rebooted)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] relax-xsm-sr-check = true (allow storage migration when SRs have been mirrored out-of-band (and have matching SR uuids))
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] include-console-username-in-error = true (Allow displaying user names in XenCenter)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] disable-logging-for = [  ] (space-separated list of modules to suppress logging from)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] disable-dbsync-for = [  ] (space-separated list of database synchronisation actions to skip)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] xenopsd-queues = org.xen.xapi.xenops.classic (list of xenopsd instances to manage)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] xenopsd-default = org.xen.xapi.xenops.classic (default xenopsd to use)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] nvidia-whitelist = /usr/share/nvidia/vgpu/vgpuConfig.xml (path to the NVidia vGPU whitelist file)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] igd-passthru-vendor-whitelist = [ 8086 ] (list of PCI vendor IDs for integrated graphics passthrough (space-separated))
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] gvt-g-whitelist = /etc/gvt-g-whitelist (path to the GVT-g whitelist file)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] gvt-g-supported = true (indicates that this server still support intel gvt_g vGPU)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] mxgpu-whitelist = /etc/mxgpu-whitelist (path to the AMD whitelist file)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] pass-through-pif-carrier = false (reflect physical interface carrier information to VMs by default)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] cluster-stack-default = xhad (Default cluster stack (HA))
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] gpumon_stop_timeout = 10. (Time to wait after attempting to stop gpumon when launching a vGPU-enabled VM.)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] reboot_required_hfxs = /run/reboot-required.hfxs (File to query hotfix uuids which require reboot)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] xen_livepatch_list = /usr/sbin/xen-livepatch list (Command to query current xen livepatch list)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] kpatch_list = /usr/sbin/kpatch list (Command to query current kernel patch list)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] modprobe_path = /usr/sbin/modprobe (Location of the modprobe(8) command: should match $(which modprobe))
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] db_idempotent_map = false (True if the add_to_<map> API calls should be idempotent)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] use-event-next = true (Use deprecated Event.next instead of Event.from)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] nvidia_multi_vgpu_enabled_driver_versions = 430.42,430.62,440.00+ (list of nvidia host driver versions with multiple vGPU supported.\x0A  if a version end with +, it means any driver version greater or equal than that version)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] nvidia_t4_sriov = default - Infer NVidia GPU addressing mode from vgpuConfig.xml (Use of SR-IOV for Nvidia GPUs; 'true', 'false', 'default'.)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] create-tools-sr = true (Indicates whether to create an SR for Tools ISOs)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] allow-host-sched-gran-modification = true (Allows to modify the host's scheduler granularity)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] use-xmlrpc = true (Use XMLRPC (deprecated) for internal communication or JSONRPC)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] extauth_ad_backend = winbind (Which AD backend used to talk to DC)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] winbind_kerberos_encryption_type = all (Encryption types to use when operating as Kerberos client [strong|legacy|all])
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] winbind_set_machine_account_kerberos_encryption_type = false (Whether set machine account encryption type (msDS-SupportedEncryptionTypes) on domain controller)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] winbind_allow_kerberos_auth_fallback = false (Whether to allow fallback to other auth on kerberos failure)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] winbind_scan_trusted_domains = false (Whether to periodically scan trusted domains)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] winbind_keep_configuration = false (Whether to clear winbind configuration when join domain failed or leave domain)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] hsts_max_age = -1 (number of seconds after the reception of the STS header field, during which the UA as a known HSTS Host (default = -1 means HSTS is disabled))
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] website-https-only = false (Allow access to the internal website using HTTPS only (no HTTP))
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] migration-https-only = true (Exclusively use HTTPS for VM migration)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] repository-domain-name-allowlist = [  ] (space-separated list of allowed domain name in base URL in repository.)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] repository-url-blocklist = [  ] (space-separated list of blocked URL patterns in base URL in repository.)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] repository-gpgcheck = true (turn gpgcheck on/off)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] repository-gpgkey-name =  (The default name of gpg key file used by YUM and RPM to verify metadata and packages in repository)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] failed-login-alert-freq = 3600 (Frequency at which we alert any failed logins (in seconds; default=3600s))
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] cert-expiration-days = 3650 (Number of days a refreshed certificate will be valid; it defaults to 10 years.)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] messages-limit = 10000 (Maximum number of messages kept before deleting oldest ones.)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] evacuation-batch-size = 10 (The number of VMs evacauted from a host in parallel.)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] ignore-vtpm-unimplemented = false (Do not raise errors on use-cases where VTPM codepaths are not finished.)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] allow-custom-uefi-certs = true (Enable (true) or Disable (false) setting a custom location for varstored UEFI certificates)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] server-cert-group-id = -1 (The group id of server ssl certificate file.)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] export-interval = 30. (The interval for exports in Tracing)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] export-chunk-size = 10000 (The span chunk size for exports in Tracing)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] max-spans = 10000 (The maximum amount of spans that can be in a trace in Tracing)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] max-traces = 10000 (The maximum number of active traces going on in Tracing)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] prefer-nbd-attach = false (Use NBD to attach disks to the control domain.)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] observer-max-file-size = 1048576 (The maximum size of log files for saving spans)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] nvidia-gpumon-detach = false (On VM start, detach the NVML library rather than stopping gpumon)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] compress-tracing-files = true (Enable compression of the tracing log files)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] observer-endpoint-http-enabled = false (Enable http endpoints to be used by observers)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] observer-endpoint-https-enabled = false (Enable https endpoints to be used by observers)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] observer-experimental-components = smapi (Comma-separated list of experimental observer components)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] disable-webserver = false (Disable the host webserver)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] tgroups-enabled = false (Turn on tgroups classification)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] event_from_delay = 0.000000,0.050000 (delays in seconds before the API call, and between internal recursive calls, separated with a comma)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] event_from_task_delay = 0.000000,0.050000 (delays in seconds before the API call, and between internal recursive calls, separated with a comma)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] event_next_delay = 0.200000,0.050000 (delays in seconds before the API call, and between internal recursive calls, separated with a comma)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] drivertool = /usr/sbin/driver-tool (Path to drivertool for selecting host driver variants)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] reuse-pool-sessions = false (Enable the reuse of pool sessions)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] validate-reusable-pool-session = false (Enable validation of reusable pool sessions before use)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] ssh-auto-mode = false (Defaults to true; overridden to false via /etc/xapi.conf.d/ssh-auto-mode.conf(e.g., in XenServer 8))
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] secure-boot-efi-path = /sys/firmware/efi/efivars/SecureBoot-8be4df61-93ca-11d2-aa0d-00e098032b8c (Path to secure boot status file)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] vm-sysprep-enabled = true (Enable VM.sysprep API)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] vm-sysprep-wait = 5. (Time in seconds to wait for VM to recognise inserted CD)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] vhd-legacy-blocks-format = false (Choose whether legacy/sparse block format will be used for determining allocated VHD clusters)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] proxy_poll_period_timeout = 5. (Timeout (in seconds) for event polling in network proxy loops. When positive, the proxy will wake up periodically to check tasks like vnc idle timeouts or perform other maintenance tasks. Set to -1 to wait indefinitely for network events without periodic wake-ups.)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] max-span-depth = 100 (The maximum depth to which spans are recorded in a trace in Tracing)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] https-only-default = false (Only expose HTTPS service, disable HTTP/80 in firewall when set to true)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] firewall-backend = iptables (Firewall backend. iptables (in XS 8) or firewalld (in XS 9 or later XS version))
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] dynamic-control-firewalld-service = true (Enable dynamic control firewalld service)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] ntp-service = chronyd (Name of the NTP service to manage)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] ntp-config-path = /etc/chrony.conf (Path to the ntp configuration file)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] ntp-dhcp-script-path = /etc/dhcp/dhclient.d/chrony.sh (Path to the ntp dhcp script file)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] ntp-dhcp-dir = /run/chrony-dhcp (Path to the ntp dhcp directory)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] ntp-client-path = /usr/bin/chronyc (Path to the ntp client binary)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] timedatectl = /usr/bin/timedatectl (Path to the timedatectl executable)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] legacy-factory-ntp-servers = [ 3.centos.pool.ntp.org; 2.centos.pool.ntp.org; 1.centos.pool.ntp.org; 0.centos.pool.ntp.org ] (space-separated list of legacy default NTP servers)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] factory-ntp-servers = [ 3.centos.pool.ntp.org; 2.centos.pool.ntp.org; 1.centos.pool.ntp.org; 0.centos.pool.ntp.org ] (space-separated list of default NTP servers)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] post-install-scripts-dir = /opt/xensource/packages/post-install-scripts (Directory containing trusted guest provisioning scripts)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] gpg-homedir = /opt/xensource/gpg (Passed as --homedir to gpg commands)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] xen-cmdline = /opt/xensource/libexec/xen-cmdline (Path to xen-cmdline binary)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] cluster-stack-root = /usr/libexec/xapi/cluster-stack (Directory containing collections of HA tools and scripts)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] web-dir = /opt/xensource/www (Directory to export fileserver)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] sm-dir = /opt/xensource/sm (Directory containing SM plugins)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] udhcpd-skel = /etc/xensource/udhcpd.skel (Skeleton config for udhcp)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] db-config-file = /etc/xensource/db.conf (Database configuration file)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] pool_config_file = /etc/xensource/pool.conf (Pool configuration file)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] udevadm = /usr/sbin/udevadm (Path to udevadm command)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] dracut = /usr/bin/dracut (Path to dracut command)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] depmod = /usr/sbin/depmod (Path to depmod command)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] rpm-cmd = /usr/bin/rpm (Path to rpm command)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] modifyrepo-cmd = /usr/bin/modifyrepo_c (Path to modifyrepo command)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] createrepo-cmd = /usr/bin/createrepo_c (Path to createrepo command)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] set-iscsi-initiator = /opt/xensource/libexec/set-iscsi-initiator (Path to set-iscsi-initiator script)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] openssl_path = /usr/bin/openssl (Path for openssl command to generate RSA keys)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] gencert = /opt/xensource/libexec/gencert (command to generate SSL certificates to be used by XAPI)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] alert-certificate-check = /opt/xensource/libexec/alert-certificate-check (Path to alert-certificate-check, which generates alerts on about-to-expire server certificates.)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] systemctl = /usr/bin/systemctl (Control the systemd system and service manager)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] list_domains = /usr/bin/list_domains (Path to the list_domains command)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] xen-cmdline-script = /opt/xensource/libexec/xen-cmdline (Path to xen-cmdline script)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] static-vdis = /opt/xensource/bin/static-vdis (Path to static-vdis script)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] xsh = /opt/xensource/bin/xsh (Path to xsh binary)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] xe-toolstack-restart = /opt/xensource/bin/xe-toolstack-restart (Path to xe-toolstack-restart script)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] xe = /usr/bin/xe (Path to xe CLI binary)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] host-restore = /opt/xensource/libexec/host-restore (Path to host-restore)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] host-backup = /opt/xensource/libexec/host-backup (Path to host-backup)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] upload-wrapper = /opt/xensource/libexec/upload-wrapper (Used by Host_crashdump.upload)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] update-mh-info = /opt/xensource/libexec/update-mh-info (Executed when changing the management interface)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] logs-download = /opt/xensource/libexec/logs-download (Used by /get_host_logs_download HTTP handler)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] xe-syslog-reconfigure = /opt/xensource/libexec/xe-syslog-reconfigure (Path to xe-syslog-reconfigure)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] set-hostname = /opt/xensource/libexec/set-hostname (Path to set-hostname)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] host-bugreport-upload = /opt/xensource/libexec/host-bugreport-upload (Path to host-bugreport-upload)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] fence = /opt/xensource/libexec/fence (Path to fence binary, used for HA host fencing)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] qcow_stream_tool = /usr/bin/qcow-stream-tool (Path to qcow-stream-tool)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] qcow_to_stdout = /opt/xensource/libexec/qcow2-to-stdout.py (Path to qcow-to-stdout script)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] vhd-tool = /usr/bin/vhd-tool (Path to vhd-tool)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] sparse_dd = /usr/libexec/xapi/sparse_dd (Path to sparse_dd)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] redo-log-block-device-io = /opt/xensource/libexec/block_device_io (Used by the redo log for block device I/O)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] pbis-force-domain-leave-script = /opt/xensource/libexec/pbis-force-domain-leave (Executed when PBIS domain-leave fails)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] busybox = /usr/bin/busybox (Swiss army knife executable - used as DHCP server)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] startup-script-hook = xapi-startup-script (Executed during startup)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] rolling-upgrade-script-hook = /opt/xensource/libexec/xapi-rolling-upgrade (Executed when a rolling upgrade is detected starting or stopping)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] xapi-message-script = /opt/xensource/libexec/mail-alarm (Executed when messages are generated if email feature is disabled)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] non-managed-pifs = /opt/xensource/libexec/bfs-interfaces (Executed during PIF.scan to find out which NICs should not be managed by xapi)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] domain_join_cli_cmd = /opt/pbis/bin/domainjoin-cli (Command to manage pbis related service)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] update-issue = /usr/sbin/update-issue (Running update-service when configuring the management interface)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] killall = /usr/bin/killall (Executed to kill process)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] nbd-firewall-config = /opt/xensource/libexec/nbd-firewall-config.sh (Executed after NBD-related networking changes to configure the firewall for NBD)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] firewall-port-config = /etc/xapi.d/plugins/firewall-port (Executed when starting/stopping xapi-clusterd to configure firewall port)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] firewall-cmd = /usr/bin/firewall-cmd (Executed when enable/disable a service on a firewalld zone)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] firewall-cmd-wrapper = /usr/bin/firewall-cmd-wrapper (Executed when enable/disable a service on a firewalld zone and interface)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] nbd_client_manager = /opt/xensource/libexec/nbd_client_manager.py (Executed to safely connect to and disconnect from NBD devices using nbd-client)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] varstore-rm = /usr/bin/varstore-rm (Executed to clear certain UEFI variables during clone)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] varstore-sb-state = /usr/bin/varstore-sb-state (Executed to edit the SecureBoot state of a VM)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] varstore-ls = /usr/bin/varstore-ls (Executed to list the UEFI variables of a VM)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] varstore_dir = /var/lib/varstored (Path to local varstored directory)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] nvidia-sriov-manage = /usr/lib/nvidia/sriov-manage (Path to NVIDIA sriov-manage script)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] gen_pool_secret_script = /usr/bin/pool_secret_wrapper (Generates new pool secrets)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] samba administration tool = /usr/bin/net (Executed to manage external auth with AD like join and leave domain)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] Samba TDB (Trivial Database) management tool = /usr/bin/tdbtool (Executed to manage Samba Database)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] winbind query tool = /usr/bin/wbinfo (Query information from winbind daemon)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] SQLite database  management tool = /usr/bin/sqlite3 (Executed to manage SQlite Database, like PBIS database)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] yum-cmd = /usr/bin/yum (Path to yum command)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] dnf-cmd = /usr/bin/dnf (Path to dnf command)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] reposync-cmd = /usr/bin/reposync (Path to reposync command)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] yum-config-manager-cmd = /usr/bin/yum-config-manager (Path to yum-config-manager command)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] c_rehash = /usr/bin/c_rehash (Path to regenerate CA store)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] fcoe-driver = /opt/xensource/libexec/fcoe_driver (Execute during PIF unplug to get the lun devices related with the ether interface of the PIF)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] pvsproxy_close_cache_vdi = /opt/citrix/pvsproxy/close-cache-vdi.sh (Path to close-cache-vdi.sh)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] genisoimage = /usr/bin/genisoimage (Path to genisoimage)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] pool_secret_path = /etc/xensource/ptoken (Pool configuration file)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] udhcpd-conf = /etc/xensource/udhcpd.conf (Optional configuration file for udchp)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] remote-db-conf-file = /etc/xensource/remote.db.conf (Where to store information about remote databases)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] logconfig = /etc/xensource/log.conf (Configure the logging policy)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] cpu-info-file = /etc/xensource/boot_time_cpus (Where to cache boot-time CPU info)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] server-cert-path = /etc/xensource/xapi-ssl.pem (Path to server ssl certificate)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] server-cert-internal-path = /etc/xensource/xapi-pool-tls.pem (Path to server certificate used for host-to-host TLS connections)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] stunnel-bundle-path = /etc/stunnel/xapi-stunnel-ca-bundle.pem (Path to stunnel trust bundle)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] pool-bundle-path = /etc/stunnel/xapi-pool-ca-bundle.pem (Path to pool trust bundle)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] iscsi_initiatorname = /etc/iscsi/initiatorname.iscsi (Path to the initiatorname.iscsi file)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] master-scripts-dir = /etc/xensource/master.d (Scripts to execute when transitioning pool role)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] packs-dir = /etc/xensource/installed-repos (Directory containing supplemental pack data)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] xapi-hooks-root = /etc/xapi.d (Root directory for xapi hooks)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] xapi-plugins-root = /etc/xapi.d/plugins (Optional directory containing XenAPI plugins)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] xapi-extensions-root = /etc/xapi.d/extensions (Optional directory containing XenAPI extensions)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] static-vdis-root = /etc/xensource/static-vdis (Optional directory for configuring static VDIs)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] tools-sr-dir = /opt/xensource/packages/iso (Directory containing tools ISO)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] trusted-pool-certs-dir = /etc/stunnel/certs-pool (Directory containing certs of trusted hosts)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] trusted-certs-dir = /etc/stunnel/certs (Directory containing certs of other trusted entities)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] trace-log-dir = /var/log/dt/zipkinv2/json (Directory for storing traces exported to logs)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] pool-recommendations-dir = /etc/xapi.pool-recommendations.d (Directory containing files with recommendations in key=value format)
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] Xcp_service.adjust_timeslice: Setting timeslice to 0.050s
> Apr 30 20:26:41 Hostname xapi: [debug||0 ||xapi] Xcp_service.adjust_timeslice: Timeslice same as or larger than OCaml's default: not setting
> Apr 30 20:26:41 Hostname xapi: [debug||0 |server_init D:97ec80851062|startup] task [XAPI SERVER STARTING]
> Apr 30 20:26:41 Hostname xapi: [debug||0 |server_init D:97ec80851062|dummytaskhelper] task XAPI SERVER STARTING D:a6bd9ee95f3f created by task D:97ec80851062
> Apr 30 20:26:41 Hostname xapi: [debug||0 |XAPI SERVER STARTING D:a6bd9ee95f3f|xapi] (Re)starting xapi, pid: 32351
> Apr 30 20:26:41 Hostname xapi: [debug||0 |XAPI SERVER STARTING D:a6bd9ee95f3f|xapi] on_system_boot=false pool_role=master
> Apr 30 20:26:41 Hostname xapi: [debug||0 |server_init D:97ec80851062|startup] task [Parsing inventory file]
> Apr 30 20:26:41 Hostname xapi: [debug||0 |server_init D:97ec80851062|dummytaskhelper] task Parsing inventory file D:60c4ea2c923a created by task D:97ec80851062
> Apr 30 20:26:41 Hostname xapi: [debug||0 |server_init D:97ec80851062|startup] task [Setting stunnel timeout]
> Apr 30 20:26:41 Hostname xapi: [debug||0 |server_init D:97ec80851062|dummytaskhelper] task Setting stunnel timeout D:b5cf96d91229 created by task D:97ec80851062
> Apr 30 20:26:41 Hostname xapi: [debug||0 |Setting stunnel timeout D:b5cf96d91229|xapi] Using default stunnel timeout (usually 43200)
> Apr 30 20:26:41 Hostname xapi: [debug||0 |server_init D:97ec80851062|startup] task [Initialising local database]
> Apr 30 20:26:41 Hostname xapi: [debug||0 |server_init D:97ec80851062|dummytaskhelper] task Initialising local database D:41bef822b055 created by task D:97ec80851062
> Apr 30 20:26:41 Hostname xapi: [debug||0 |Initialising local database D:41bef822b055|hashtbl_xml] Converting dtd
> Apr 30 20:26:41 Hostname xapi: [debug||0 |Initialising local database D:41bef822b055|localdb] host_auto_enable = true
> Apr 30 20:26:41 Hostname xapi: [debug||0 |Initialising local database D:41bef822b055|localdb] host_restarted_cleanly = false
> Apr 30 20:26:41 Hostname xapi: [debug||0 |Initialising local database D:41bef822b055|localdb] master_scripts = false
> Apr 30 20:26:41 Hostname xapi: [debug||0 |Initialising local database D:41bef822b055|localdb] ha.armed = false
> Apr 30 20:26:41 Hostname xapi: [debug||0 |Initialising local database D:41bef822b055|localdb] host_disabled_until_reboot = false
> Apr 30 20:26:41 Hostname xapi: [debug||0 |server_init D:97ec80851062|startup] task [Loading DHCP leases]
> Apr 30 20:26:41 Hostname xapi: [debug||0 |server_init D:97ec80851062|dummytaskhelper] task Loading DHCP leases D:e0ec2dc65f36 created by task D:97ec80851062
> Apr 30 20:26:41 Hostname xapi: [ info||0 |Loading DHCP leases D:e0ec2dc65f36|xapi_udhcpd] Caught exception Unix.Unix_error(Unix.ENOENT, "open", "/var/lib/xcp/dhcp-leases.db") loading /var/lib/xcp/dhcp-leases.db: creating new empty leases database
> Apr 30 20:26:41 Hostname xapi: [debug||0 |server_init D:97ec80851062|startup] task [Reading pool secret]
> Apr 30 20:26:41 Hostname xapi: [debug||0 |server_init D:97ec80851062|dummytaskhelper] task Reading pool secret D:97f682043e63 created by task D:97ec80851062
> Apr 30 20:26:41 Hostname xapi: [debug||0 |server_init D:97ec80851062|startup] task [Logging xapi version info]
> Apr 30 20:26:41 Hostname xapi: [debug||0 |server_init D:97ec80851062|dummytaskhelper] task Logging xapi version info D:9a77de11bb6d created by task D:97ec80851062
> Apr 30 20:26:41 Hostname xapi: [debug||0 |Logging xapi version info D:9a77de11bb6d|Xapi_config] Server configuration:
> Apr 30 20:26:41 Hostname xapi: [debug||0 |Logging xapi version info D:9a77de11bb6d|Xapi_config] product_version: 8.3.0
> Apr 30 20:26:41 Hostname xapi: [debug||0 |Logging xapi version info D:9a77de11bb6d|Xapi_config] product_brand: XCP-ng
> Apr 30 20:26:41 Hostname xapi: [debug||0 |Logging xapi version info D:9a77de11bb6d|Xapi_config] platform_version: 3.4.0
> Apr 30 20:26:41 Hostname xapi: [debug||0 |Logging xapi version info D:9a77de11bb6d|Xapi_config] platform_name: XCP
> Apr 30 20:26:41 Hostname xapi: [debug||0 |Logging xapi version info D:9a77de11bb6d|Xapi_config] build_number: 8.3.0
> Apr 30 20:26:41 Hostname xapi: [debug||0 |Logging xapi version info D:9a77de11bb6d|Xapi_config] version: 26.1.3
> Apr 30 20:26:41 Hostname xapi: [debug||0 |server_init D:97ec80851062|startup] task [Setting signal handlers]
> Apr 30 20:26:41 Hostname xapi: [debug||0 |server_init D:97ec80851062|dummytaskhelper] task Setting signal handlers D:56d827b002cc created by task D:97ec80851062
> Apr 30 20:26:41 Hostname xapi: [debug||0 |server_init D:97ec80851062|startup] task [Initialising random number generator]
> Apr 30 20:26:41 Hostname xapi: [debug||0 |server_init D:97ec80851062|dummytaskhelper] task Initialising random number generator D:4ed9d35dca03 created by task D:97ec80851062
> Apr 30 20:26:41 Hostname xapi: [debug||0 |server_init D:97ec80851062|startup] task [Initialise TLS verification]
> Apr 30 20:26:41 Hostname xapi: [debug||0 |server_init D:97ec80851062|dummytaskhelper] task Initialise TLS verification D:c2391d2954a7 created by task D:97ec80851062
> Apr 30 20:26:41 Hostname xapi: [ info||0 |Initialise TLS verification D:c2391d2954a7|xapi] TLS verification is enabled: /var/xapi/verify-certificates is present
> Apr 30 20:26:41 Hostname xapi: [ info||0 |Initialise TLS verification D:c2391d2954a7|Stunnel_client] enabling default tls verification
> Apr 30 20:26:41 Hostname xapi: [debug||0 |server_init D:97ec80851062|startup] task [Running startup check]
> Apr 30 20:26:41 Hostname xapi: [debug||0 |server_init D:97ec80851062|dummytaskhelper] task Running startup check D:01fbc4364e61 created by task D:97ec80851062
> Apr 30 20:26:41 Hostname xapi: [debug||0 |Running startup check D:01fbc4364e61|sanitycheck] Binary appears to be correctly linked
> Apr 30 20:26:41 Hostname xapi: [debug||0 |server_init D:97ec80851062|startup] task [Initialize cgroups via tgroup]
> Apr 30 20:26:41 Hostname xapi: [debug||0 |server_init D:97ec80851062|dummytaskhelper] task Initialize cgroups via tgroup D:122ba98ba402 created by task D:97ec80851062
> Apr 30 20:26:41 Hostname xapi: [debug||0 |server_init D:97ec80851062|startup] task [Registering SMAPIv1 plugins]
> Apr 30 20:26:41 Hostname xapi: [debug||0 |server_init D:97ec80851062|dummytaskhelper] task Registering SMAPIv1 plugins D:2af8589e72af created by task D:97ec80851062
> Apr 30 20:26:41 Hostname xapi: [debug||0 |Registering SMAPIv1 plugins D:2af8589e72af|sm_exec] Scanning directory /opt/xensource/sm for SM plugins
> Apr 30 20:26:41 Hostname xapi: [ warn||0 |sm_exec D:dda3f0c5514f|helpers] The database has not fully come up yet, so localhost is missing
> Apr 30 20:26:41 Hostname xapi: [error||0 |check if component smapi is enabled  D:206ed8a2473a|backtrace] Raised Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:41 Hostname xapi: [error||0 |check if component smapi is enabled  D:206ed8a2473a|backtrace] 1/4 xapi Raised at file ocaml/database/db_cache_types.ml, line 298
> Apr 30 20:26:41 Hostname xapi: [error||0 |check if component smapi is enabled  D:206ed8a2473a|backtrace] 2/4 xapi Called from file ocaml/database/db_cache_impl.ml, line 314
> Apr 30 20:26:41 Hostname xapi: [error||0 |check if component smapi is enabled  D:206ed8a2473a|backtrace] 3/4 xapi Called from file ocaml/xapi/db_actions.ml, line 22998
> Apr 30 20:26:41 Hostname xapi: [error||0 |check if component smapi is enabled  D:206ed8a2473a|backtrace] 4/4 xapi Called from file ocaml/xapi/xapi_observer_components.ml, line 83
> Apr 30 20:26:41 Hostname xapi: [error||0 |check if component smapi is enabled  D:206ed8a2473a|backtrace]
> Apr 30 20:26:41 Hostname xapi: [ warn||0 |check if component smapi is enabled  D:206ed8a2473a|xapi_observer_components] is_component_enabled(smapi) inner got exception: Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:41 Hostname xapi: [ warn||0 |sm_exec D:5d79b53eea2b|helpers] The database has not fully come up yet, so localhost is missing
> Apr 30 20:26:41 Hostname xapi: [error||0 |check if component smapi is enabled  D:c39cdb372b10|backtrace] Raised Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:41 Hostname xapi: [error||0 |check if component smapi is enabled  D:c39cdb372b10|backtrace] 1/4 xapi Raised at file ocaml/database/db_cache_types.ml, line 298
> Apr 30 20:26:41 Hostname xapi: [error||0 |check if component smapi is enabled  D:c39cdb372b10|backtrace] 2/4 xapi Called from file ocaml/database/db_cache_impl.ml, line 314
> Apr 30 20:26:41 Hostname xapi: [error||0 |check if component smapi is enabled  D:c39cdb372b10|backtrace] 3/4 xapi Called from file ocaml/xapi/db_actions.ml, line 22998
> Apr 30 20:26:41 Hostname xapi: [error||0 |check if component smapi is enabled  D:c39cdb372b10|backtrace] 4/4 xapi Called from file ocaml/xapi/xapi_observer_components.ml, line 83
> Apr 30 20:26:41 Hostname xapi: [error||0 |check if component smapi is enabled  D:c39cdb372b10|backtrace]
> Apr 30 20:26:41 Hostname xapi: [ warn||0 |check if component smapi is enabled  D:c39cdb372b10|xapi_observer_components] is_component_enabled(smapi) inner got exception: Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:42 Hostname xapi: [ warn||0 |Registering SMAPIv1 plugins D:2af8589e72af|smint] SM.feature: unknown feature ATOMIC_PAUSE
> Apr 30 20:26:42 Hostname xapi: [ warn||0 |sm_exec D:80e3a23a4e34|helpers] The database has not fully come up yet, so localhost is missing
> Apr 30 20:26:42 Hostname xapi: [error||0 |check if component smapi is enabled  D:f8903bab9b24|backtrace] Raised Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:42 Hostname xapi: [error||0 |check if component smapi is enabled  D:f8903bab9b24|backtrace] 1/4 xapi Raised at file ocaml/database/db_cache_types.ml, line 298
> Apr 30 20:26:42 Hostname xapi: [error||0 |check if component smapi is enabled  D:f8903bab9b24|backtrace] 2/4 xapi Called from file ocaml/database/db_cache_impl.ml, line 314
> Apr 30 20:26:42 Hostname xapi: [error||0 |check if component smapi is enabled  D:f8903bab9b24|backtrace] 3/4 xapi Called from file ocaml/xapi/db_actions.ml, line 22998
> Apr 30 20:26:42 Hostname xapi: [error||0 |check if component smapi is enabled  D:f8903bab9b24|backtrace] 4/4 xapi Called from file ocaml/xapi/xapi_observer_components.ml, line 83
> Apr 30 20:26:42 Hostname xapi: [error||0 |check if component smapi is enabled  D:f8903bab9b24|backtrace]
> Apr 30 20:26:42 Hostname xapi: [ warn||0 |check if component smapi is enabled  D:f8903bab9b24|xapi_observer_components] is_component_enabled(smapi) inner got exception: Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:42 Hostname xapi: [ warn||0 |Registering SMAPIv1 plugins D:2af8589e72af|smint] SM.feature: unknown feature ATOMIC_PAUSE
> Apr 30 20:26:42 Hostname xapi: [ warn||0 |sm_exec D:fee34d50cd5f|helpers] The database has not fully come up yet, so localhost is missing
> Apr 30 20:26:42 Hostname xapi: [error||0 |check if component smapi is enabled  D:10e11cca8d02|backtrace] Raised Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:42 Hostname xapi: [error||0 |check if component smapi is enabled  D:10e11cca8d02|backtrace] 1/4 xapi Raised at file ocaml/database/db_cache_types.ml, line 298
> Apr 30 20:26:42 Hostname xapi: [error||0 |check if component smapi is enabled  D:10e11cca8d02|backtrace] 2/4 xapi Called from file ocaml/database/db_cache_impl.ml, line 314
> Apr 30 20:26:42 Hostname xapi: [error||0 |check if component smapi is enabled  D:10e11cca8d02|backtrace] 3/4 xapi Called from file ocaml/xapi/db_actions.ml, line 22998
> Apr 30 20:26:42 Hostname xapi: [error||0 |check if component smapi is enabled  D:10e11cca8d02|backtrace] 4/4 xapi Called from file ocaml/xapi/xapi_observer_components.ml, line 83
> Apr 30 20:26:42 Hostname xapi: [error||0 |check if component smapi is enabled  D:10e11cca8d02|backtrace]
> Apr 30 20:26:42 Hostname xapi: [ warn||0 |check if component smapi is enabled  D:10e11cca8d02|xapi_observer_components] is_component_enabled(smapi) inner got exception: Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:43 Hostname xapi: [ warn||0 |Registering SMAPIv1 plugins D:2af8589e72af|smint] SM.feature: unknown feature ATOMIC_PAUSE
> Apr 30 20:26:43 Hostname xapi: [ warn||0 |sm_exec D:4bee50fb2560|helpers] The database has not fully come up yet, so localhost is missing
> Apr 30 20:26:43 Hostname xapi: [error||0 |check if component smapi is enabled  D:69ac0ef3dc50|backtrace] Raised Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:43 Hostname xapi: [error||0 |check if component smapi is enabled  D:69ac0ef3dc50|backtrace] 1/4 xapi Raised at file ocaml/database/db_cache_types.ml, line 298
> Apr 30 20:26:43 Hostname xapi: [error||0 |check if component smapi is enabled  D:69ac0ef3dc50|backtrace] 2/4 xapi Called from file ocaml/database/db_cache_impl.ml, line 314
> Apr 30 20:26:43 Hostname xapi: [error||0 |check if component smapi is enabled  D:69ac0ef3dc50|backtrace] 3/4 xapi Called from file ocaml/xapi/db_actions.ml, line 22998
> Apr 30 20:26:43 Hostname xapi: [error||0 |check if component smapi is enabled  D:69ac0ef3dc50|backtrace] 4/4 xapi Called from file ocaml/xapi/xapi_observer_components.ml, line 83
> Apr 30 20:26:43 Hostname xapi: [error||0 |check if component smapi is enabled  D:69ac0ef3dc50|backtrace]
> Apr 30 20:26:43 Hostname xapi: [ warn||0 |check if component smapi is enabled  D:69ac0ef3dc50|xapi_observer_components] is_component_enabled(smapi) inner got exception: Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:43 Hostname xapi: [ warn||0 |sm_exec D:d881d27e8418|helpers] The database has not fully come up yet, so localhost is missing
> Apr 30 20:26:43 Hostname xapi: [error||0 |check if component smapi is enabled  D:f9824fe54110|backtrace] Raised Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:43 Hostname xapi: [error||0 |check if component smapi is enabled  D:f9824fe54110|backtrace] 1/4 xapi Raised at file ocaml/database/db_cache_types.ml, line 298
> Apr 30 20:26:43 Hostname xapi: [error||0 |check if component smapi is enabled  D:f9824fe54110|backtrace] 2/4 xapi Called from file ocaml/database/db_cache_impl.ml, line 314
> Apr 30 20:26:43 Hostname xapi: [error||0 |check if component smapi is enabled  D:f9824fe54110|backtrace] 3/4 xapi Called from file ocaml/xapi/db_actions.ml, line 22998
> Apr 30 20:26:43 Hostname xapi: [error||0 |check if component smapi is enabled  D:f9824fe54110|backtrace] 4/4 xapi Called from file ocaml/xapi/xapi_observer_components.ml, line 83
> Apr 30 20:26:43 Hostname xapi: [error||0 |check if component smapi is enabled  D:f9824fe54110|backtrace]
> Apr 30 20:26:43 Hostname xapi: [ warn||0 |check if component smapi is enabled  D:f9824fe54110|xapi_observer_components] is_component_enabled(smapi) inner got exception: Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:43 Hostname xapi: [ warn||0 |Registering SMAPIv1 plugins D:2af8589e72af|smint] SM.feature: unknown feature ATOMIC_PAUSE
> Apr 30 20:26:43 Hostname xapi: [ warn||0 |sm_exec D:570f3e1bcda7|helpers] The database has not fully come up yet, so localhost is missing
> Apr 30 20:26:43 Hostname xapi: [error||0 |check if component smapi is enabled  D:575aa4600c78|backtrace] Raised Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:43 Hostname xapi: [error||0 |check if component smapi is enabled  D:575aa4600c78|backtrace] 1/4 xapi Raised at file ocaml/database/db_cache_types.ml, line 298
> Apr 30 20:26:43 Hostname xapi: [error||0 |check if component smapi is enabled  D:575aa4600c78|backtrace] 2/4 xapi Called from file ocaml/database/db_cache_impl.ml, line 314
> Apr 30 20:26:43 Hostname xapi: [error||0 |check if component smapi is enabled  D:575aa4600c78|backtrace] 3/4 xapi Called from file ocaml/xapi/db_actions.ml, line 22998
> Apr 30 20:26:43 Hostname xapi: [error||0 |check if component smapi is enabled  D:575aa4600c78|backtrace] 4/4 xapi Called from file ocaml/xapi/xapi_observer_components.ml, line 83
> Apr 30 20:26:43 Hostname xapi: [error||0 |check if component smapi is enabled  D:575aa4600c78|backtrace]
> Apr 30 20:26:43 Hostname xapi: [ warn||0 |check if component smapi is enabled  D:575aa4600c78|xapi_observer_components] is_component_enabled(smapi) inner got exception: Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:44 Hostname xapi: [ warn||0 |Registering SMAPIv1 plugins D:2af8589e72af|smint] SM.feature: unknown feature ATOMIC_PAUSE
> Apr 30 20:26:44 Hostname xapi: [ warn||0 |sm_exec D:e03c72296895|helpers] The database has not fully come up yet, so localhost is missing
> Apr 30 20:26:44 Hostname xapi: [error||0 |check if component smapi is enabled  D:edb2a0ce6ac3|backtrace] Raised Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:44 Hostname xapi: [error||0 |check if component smapi is enabled  D:edb2a0ce6ac3|backtrace] 1/4 xapi Raised at file ocaml/database/db_cache_types.ml, line 298
> Apr 30 20:26:44 Hostname xapi: [error||0 |check if component smapi is enabled  D:edb2a0ce6ac3|backtrace] 2/4 xapi Called from file ocaml/database/db_cache_impl.ml, line 314
> Apr 30 20:26:44 Hostname xapi: [error||0 |check if component smapi is enabled  D:edb2a0ce6ac3|backtrace] 3/4 xapi Called from file ocaml/xapi/db_actions.ml, line 22998
> Apr 30 20:26:44 Hostname xapi: [error||0 |check if component smapi is enabled  D:edb2a0ce6ac3|backtrace] 4/4 xapi Called from file ocaml/xapi/xapi_observer_components.ml, line 83
> Apr 30 20:26:44 Hostname xapi: [error||0 |check if component smapi is enabled  D:edb2a0ce6ac3|backtrace]
> Apr 30 20:26:44 Hostname xapi: [ warn||0 |check if component smapi is enabled  D:edb2a0ce6ac3|xapi_observer_components] is_component_enabled(smapi) inner got exception: Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:44 Hostname xapi: [ warn||0 |sm_exec D:b3cd99521e5e|helpers] The database has not fully come up yet, so localhost is missing
> Apr 30 20:26:44 Hostname xapi: [error||0 |check if component smapi is enabled  D:d11624b2cfa9|backtrace] Raised Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:44 Hostname xapi: [error||0 |check if component smapi is enabled  D:d11624b2cfa9|backtrace] 1/4 xapi Raised at file ocaml/database/db_cache_types.ml, line 298
> Apr 30 20:26:44 Hostname xapi: [error||0 |check if component smapi is enabled  D:d11624b2cfa9|backtrace] 2/4 xapi Called from file ocaml/database/db_cache_impl.ml, line 314
> Apr 30 20:26:44 Hostname xapi: [error||0 |check if component smapi is enabled  D:d11624b2cfa9|backtrace] 3/4 xapi Called from file ocaml/xapi/db_actions.ml, line 22998
> Apr 30 20:26:44 Hostname xapi: [error||0 |check if component smapi is enabled  D:d11624b2cfa9|backtrace] 4/4 xapi Called from file ocaml/xapi/xapi_observer_components.ml, line 83
> Apr 30 20:26:44 Hostname xapi: [error||0 |check if component smapi is enabled  D:d11624b2cfa9|backtrace]
> Apr 30 20:26:44 Hostname xapi: [ warn||0 |check if component smapi is enabled  D:d11624b2cfa9|xapi_observer_components] is_component_enabled(smapi) inner got exception: Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:45 Hostname xapi: [ warn||0 |Registering SMAPIv1 plugins D:2af8589e72af|smint] SM.feature: unknown feature ATOMIC_PAUSE
> Apr 30 20:26:45 Hostname xapi: [ warn||0 |sm_exec D:1f5d77c80117|helpers] The database has not fully come up yet, so localhost is missing
> Apr 30 20:26:45 Hostname xapi: [error||0 |check if component smapi is enabled  D:1206246df06c|backtrace] Raised Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:45 Hostname xapi: [error||0 |check if component smapi is enabled  D:1206246df06c|backtrace] 1/4 xapi Raised at file ocaml/database/db_cache_types.ml, line 298
> Apr 30 20:26:45 Hostname xapi: [error||0 |check if component smapi is enabled  D:1206246df06c|backtrace] 2/4 xapi Called from file ocaml/database/db_cache_impl.ml, line 314
> Apr 30 20:26:45 Hostname xapi: [error||0 |check if component smapi is enabled  D:1206246df06c|backtrace] 3/4 xapi Called from file ocaml/xapi/db_actions.ml, line 22998
> Apr 30 20:26:45 Hostname xapi: [error||0 |check if component smapi is enabled  D:1206246df06c|backtrace] 4/4 xapi Called from file ocaml/xapi/xapi_observer_components.ml, line 83
> Apr 30 20:26:45 Hostname xapi: [error||0 |check if component smapi is enabled  D:1206246df06c|backtrace]
> Apr 30 20:26:45 Hostname xapi: [ warn||0 |check if component smapi is enabled  D:1206246df06c|xapi_observer_components] is_component_enabled(smapi) inner got exception: Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:45 Hostname xapi: [ warn||0 |Registering SMAPIv1 plugins D:2af8589e72af|smint] SM.feature: unknown feature ATOMIC_PAUSE
> Apr 30 20:26:45 Hostname xapi: [ warn||0 |sm_exec D:7836fe923481|helpers] The database has not fully come up yet, so localhost is missing
> Apr 30 20:26:45 Hostname xapi: [error||0 |check if component smapi is enabled  D:cd24b3ab5a16|backtrace] Raised Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:45 Hostname xapi: [error||0 |check if component smapi is enabled  D:cd24b3ab5a16|backtrace] 1/4 xapi Raised at file ocaml/database/db_cache_types.ml, line 298
> Apr 30 20:26:45 Hostname xapi: [error||0 |check if component smapi is enabled  D:cd24b3ab5a16|backtrace] 2/4 xapi Called from file ocaml/database/db_cache_impl.ml, line 314
> Apr 30 20:26:45 Hostname xapi: [error||0 |check if component smapi is enabled  D:cd24b3ab5a16|backtrace] 3/4 xapi Called from file ocaml/xapi/db_actions.ml, line 22998
> Apr 30 20:26:45 Hostname xapi: [error||0 |check if component smapi is enabled  D:cd24b3ab5a16|backtrace] 4/4 xapi Called from file ocaml/xapi/xapi_observer_components.ml, line 83
> Apr 30 20:26:45 Hostname xapi: [error||0 |check if component smapi is enabled  D:cd24b3ab5a16|backtrace]
> Apr 30 20:26:45 Hostname xapi: [ warn||0 |check if component smapi is enabled  D:cd24b3ab5a16|xapi_observer_components] is_component_enabled(smapi) inner got exception: Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:46 Hostname xapi: [ warn||0 |Registering SMAPIv1 plugins D:2af8589e72af|smint] SM.feature: unknown feature ATOMIC_PAUSE
> Apr 30 20:26:46 Hostname xapi: [ warn||0 |sm_exec D:bc0ae98c6e7e|helpers] The database has not fully come up yet, so localhost is missing
> Apr 30 20:26:46 Hostname xapi: [error||0 |check if component smapi is enabled  D:b9cdc00587e5|backtrace] Raised Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:46 Hostname xapi: [error||0 |check if component smapi is enabled  D:b9cdc00587e5|backtrace] 1/4 xapi Raised at file ocaml/database/db_cache_types.ml, line 298
> Apr 30 20:26:46 Hostname xapi: [error||0 |check if component smapi is enabled  D:b9cdc00587e5|backtrace] 2/4 xapi Called from file ocaml/database/db_cache_impl.ml, line 314
> Apr 30 20:26:46 Hostname xapi: [error||0 |check if component smapi is enabled  D:b9cdc00587e5|backtrace] 3/4 xapi Called from file ocaml/xapi/db_actions.ml, line 22998
> Apr 30 20:26:46 Hostname xapi: [error||0 |check if component smapi is enabled  D:b9cdc00587e5|backtrace] 4/4 xapi Called from file ocaml/xapi/xapi_observer_components.ml, line 83
> Apr 30 20:26:46 Hostname xapi: [error||0 |check if component smapi is enabled  D:b9cdc00587e5|backtrace]
> Apr 30 20:26:46 Hostname xapi: [ warn||0 |check if component smapi is enabled  D:b9cdc00587e5|xapi_observer_components] is_component_enabled(smapi) inner got exception: Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:46 Hostname xapi: [ warn||0 |Registering SMAPIv1 plugins D:2af8589e72af|smint] SM.feature: unknown feature ATOMIC_PAUSE
> Apr 30 20:26:46 Hostname xapi: [ warn||0 |sm_exec D:864d21546e2b|helpers] The database has not fully come up yet, so localhost is missing
> Apr 30 20:26:46 Hostname xapi: [error||0 |check if component smapi is enabled  D:68a9e068b22c|backtrace] Raised Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:46 Hostname xapi: [error||0 |check if component smapi is enabled  D:68a9e068b22c|backtrace] 1/4 xapi Raised at file ocaml/database/db_cache_types.ml, line 298
> Apr 30 20:26:46 Hostname xapi: [error||0 |check if component smapi is enabled  D:68a9e068b22c|backtrace] 2/4 xapi Called from file ocaml/database/db_cache_impl.ml, line 314
> Apr 30 20:26:46 Hostname xapi: [error||0 |check if component smapi is enabled  D:68a9e068b22c|backtrace] 3/4 xapi Called from file ocaml/xapi/db_actions.ml, line 22998
> Apr 30 20:26:46 Hostname xapi: [error||0 |check if component smapi is enabled  D:68a9e068b22c|backtrace] 4/4 xapi Called from file ocaml/xapi/xapi_observer_components.ml, line 83
> Apr 30 20:26:46 Hostname xapi: [error||0 |check if component smapi is enabled  D:68a9e068b22c|backtrace]
> Apr 30 20:26:46 Hostname xapi: [ warn||0 |check if component smapi is enabled  D:68a9e068b22c|xapi_observer_components] is_component_enabled(smapi) inner got exception: Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:46 Hostname xapi: [ warn||0 |Registering SMAPIv1 plugins D:2af8589e72af|smint] SM.feature: unknown feature ATOMIC_PAUSE
> Apr 30 20:26:46 Hostname xapi: [ info||0 |Registering SMAPIv1 plugins D:2af8589e72af|sm_exec] Skipping SMAPIv1 plugin MooseFS: not in sm-plugins whitelist in configuration file
> Apr 30 20:26:46 Hostname xapi: [ info||0 |Registering SMAPIv1 plugins D:2af8589e72af|sm_exec] Skipping SMAPIv1 plugin Linstor: not in sm-plugins whitelist in configuration file
> Apr 30 20:26:46 Hostname xapi: [ warn||0 |sm_exec D:51e1232ec48f|helpers] The database has not fully come up yet, so localhost is missing
> Apr 30 20:26:46 Hostname xapi: [error||0 |check if component smapi is enabled  D:d305e8aff3a8|backtrace] Raised Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:46 Hostname xapi: [error||0 |check if component smapi is enabled  D:d305e8aff3a8|backtrace] 1/4 xapi Raised at file ocaml/database/db_cache_types.ml, line 298
> Apr 30 20:26:46 Hostname xapi: [error||0 |check if component smapi is enabled  D:d305e8aff3a8|backtrace] 2/4 xapi Called from file ocaml/database/db_cache_impl.ml, line 314
> Apr 30 20:26:46 Hostname xapi: [error||0 |check if component smapi is enabled  D:d305e8aff3a8|backtrace] 3/4 xapi Called from file ocaml/xapi/db_actions.ml, line 22998
> Apr 30 20:26:46 Hostname xapi: [error||0 |check if component smapi is enabled  D:d305e8aff3a8|backtrace] 4/4 xapi Called from file ocaml/xapi/xapi_observer_components.ml, line 83
> Apr 30 20:26:46 Hostname xapi: [error||0 |check if component smapi is enabled  D:d305e8aff3a8|backtrace]
> Apr 30 20:26:46 Hostname xapi: [ warn||0 |check if component smapi is enabled  D:d305e8aff3a8|xapi_observer_components] is_component_enabled(smapi) inner got exception: Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:47 Hostname xapi: [ warn||0 |Registering SMAPIv1 plugins D:2af8589e72af|smint] SM.feature: unknown feature ATOMIC_PAUSE
> Apr 30 20:26:47 Hostname xapi: [ warn||0 |sm_exec D:68fa33ac8747|helpers] The database has not fully come up yet, so localhost is missing
> Apr 30 20:26:47 Hostname xapi: [error||0 |check if component smapi is enabled  D:cbc9ca4660b8|backtrace] Raised Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:47 Hostname xapi: [error||0 |check if component smapi is enabled  D:cbc9ca4660b8|backtrace] 1/4 xapi Raised at file ocaml/database/db_cache_types.ml, line 298
> Apr 30 20:26:47 Hostname xapi: [error||0 |check if component smapi is enabled  D:cbc9ca4660b8|backtrace] 2/4 xapi Called from file ocaml/database/db_cache_impl.ml, line 314
> Apr 30 20:26:47 Hostname xapi: [error||0 |check if component smapi is enabled  D:cbc9ca4660b8|backtrace] 3/4 xapi Called from file ocaml/xapi/db_actions.ml, line 22998
> Apr 30 20:26:47 Hostname xapi: [error||0 |check if component smapi is enabled  D:cbc9ca4660b8|backtrace] 4/4 xapi Called from file ocaml/xapi/xapi_observer_components.ml, line 83
> Apr 30 20:26:47 Hostname xapi: [error||0 |check if component smapi is enabled  D:cbc9ca4660b8|backtrace]
> Apr 30 20:26:47 Hostname xapi: [ warn||0 |check if component smapi is enabled  D:cbc9ca4660b8|xapi_observer_components] is_component_enabled(smapi) inner got exception: Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:47 Hostname xapi: [ warn||0 |sm_exec D:d2f215cb9f9b|helpers] The database has not fully come up yet, so localhost is missing
> Apr 30 20:26:47 Hostname xapi: [error||0 |check if component smapi is enabled  D:d9772ce581f6|backtrace] Raised Db_exn.DBCache_NotFound("missing table", "Observer", "")
> Apr 30 20:26:47 Hostname xapi: [error||0 |check if component smapi is enabled  D:d9772ce581f6|backtrace] 1/4 xapi Raised at file ocaml/database/db_cache_types.ml, line 298
> Apr 30 20:26:47 Hostname xapi: [error||0 |check if component smapi is enabled  D:d9772ce581f6|backtrace] 2/4 xapi Called from file ocaml/database/db_cache_impl.ml, line 314
> Apr 30 20:26:47 Hostname xapi: [error||0 |check if component smapi is enabled  D:d9772ce581f6|backtrace] 3/4 xapi Called from file ocaml/xapi/db_actions.ml, line 22998
> Apr 30 20:26:47 Hostname xapi: [error||0 |check if component smapi is enabled  D:d9772ce581f6|backtrace] 4/4 xapi Called from file ocaml/xapi/xapi_observer_components.ml, line 83
> Apr 30 20:26:47 Hostname xapi: [error||0 |check if component smapi is enabled  D:d9772ce581f6|backtrace]
> Apr 30 20:26:47 Hostname xapi: [ warn||0 |check if component smapi is enabled  D:d9772ce581f6|xapi_observer_components] is_component_enabled(smapi) inner got exception: Db_exn.DBCache_NotFound("missing table", "Observer", "")
> ```

state.db is borked
> [!CMD]- `# xmllint --noout /var/lib/xcp/state.db`
> ```
> /var/lib/xcp/state.db:1: parser error : Start tag expected, '<' not found
> (b);b"on:end",se.isMatchIgnored&&(D=!1)}if(D){for(;b.endsParent&&b.paren
> ^
> ```

`# tail state.db -n 25` [[attachments/tail state.db.corrupt -n 25.out]]
the state db should be XML. this is not XML...

# Potential Cause
The server lost power and did not gracefully shutdown

# My Solution
```
systemctl stop xapi
mv state.db state.db.corrupt
mv state.db.generation state.db.generation.corrupt
xe-toolchain-restart
```
New database created, nothing populated

Now do an Emergency Network Reset and reboot
This will get networking on the host working

Now fix storage
```
xe sr-introduce uuid=11e549c7-968d-4148-ba56-c4c68acdd168 type=ext name-label="DATA (recovered)" content-type=user
pbd-create host-uuid=29068b64-8bf5-4924-9e48-efebdd187af7 sr-uuid=11e549c7-968d-4148-ba56-c4c68acdd168 device-config:device=/dev/sda3
xe pbd-plug uuid=d79d4ed9-ac70-d36d-7ee4-4321ce2ee572
xe sr-list
xe sr-scan uuid=11e549c7-968d-4148-ba56-c4c68acdd168 
```

Now recreate a VM with the existing disks (Hopefully your XO)
```
yum reinstall "guest-templates-json*"
xe vdi-list sr-uuid=11e549c7-968d-4148-ba56-c4c68acdd168
xe vm-install template=Other\ install\ media new-name-label='recovered1'
xe vbd-create vm-uuid=3e78c985-8839-a93b-acf7-955ff66261de vdi-uuid=5dbcbd32-c67a-44f5-8c1c-af34d53b18cd device=0 bootable=true mode=RW type=Disk
xe vif-create vm-uuid=3e78c985-8839-a93b-acf7-955ff66261de network-uuid=9efa6d67-1eb0-8364-3dfe-a6711cdceed9 device=0
xe vm-start uuid=3e78c985-8839-a93b-acf7-955ff66261de
```
Now boot it to see what it is. In my case I was right and this was my XO. Shut it down once you've identified it so you can correct the metadata
```
xe vm-param-set uuid=3e78c985-8839-a93b-acf7-955ff66261de name-label='XO-CE'
xe vm-param-set uuid=3e78c985-8839-a93b-acf7-955ff66261de memory-static-max=8589934592 memory-dynamic-max=8589934592 memory-dynamic-min=2147483648 VCPUs-max=8
```

now backup this state of the databases
```
cp state.db state.db.bak
cp state.db.generation state.db.generation.bak
cp networkd.db networkd.db.bak
cp local.db local.db.bak
```