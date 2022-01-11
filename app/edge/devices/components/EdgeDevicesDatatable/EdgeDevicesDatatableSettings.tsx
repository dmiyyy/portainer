import { TableSettingsMenuAutoRefresh } from '@/portainer/datatables/components/TableSettingsMenuAutoRefresh';
import { useTableSettings } from '@/portainer/datatables/components/useTableSettings';

import { EdgeDeviceTableSettings } from '@/edge/devices/types';

export function EdgeDevicesDatatableSettings() {
  const { settings, setTableSettings } =
    useTableSettings<EdgeDeviceTableSettings>();

  return (
    <TableSettingsMenuAutoRefresh
      value={settings.autoRefreshRate}
      onChange={handleRefreshRateChange}
    />
  );

  function handleRefreshRateChange(autoRefreshRate: number) {
    setTableSettings({ autoRefreshRate });
  }
}
