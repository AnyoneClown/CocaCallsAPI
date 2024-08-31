import { useMemo, useEffect, ReactElement } from 'react';
import { Stack, Avatar, Tooltip, Typography, CircularProgress } from '@mui/material';
import {
  GridApi,
  DataGrid,
  GridSlots,
  GridColDef,
  useGridApiRef,
  GridActionsCellItem,
  GridRenderCellParams,
  GridTreeNodeWithRender,
} from '@mui/x-data-grid';
import { rows } from 'data/customer-data';
import { stringAvatar } from 'helpers/string-avatar';
import IconifyIcon from 'components/base/IconifyIcon';
import { currencyFormat } from 'helpers/format-functions';
import CustomPagination from 'components/common/CustomPagination';
import CustomNoResultsOverlay from 'components/common/CustomNoResultsOverlay';

import { format } from 'date-fns';

const columns: GridColDef<any>[] = [
  {
    field: 'id',
    headerName: 'ID',
    resizable: false,
    minWidth: 60,
  },
  {
    field: 'email',
    headerName: 'Email',
    resizable: false,
    flex: 1,
    minWidth: 145,
  },
  {
    field: 'picture',
    headerName: 'Picture',
    renderCell: (params: GridRenderCellParams<any, any, any, GridTreeNodeWithRender>) => {
      return (
        <Tooltip title={params.row.picture} placement="top" arrow>
          <Avatar src={params.row.picture} />
        </Tooltip>
      );
    },
    resizable: false,
    flex: 1,
    minWidth: 100,
  },
  {
    field: 'is-admin',
    headerName: 'Is Admin',
    resizable: false,
    flex: 1,
    minWidth: 100,
    valueFormatter: (params) => (params.value ? 'Yes' : 'No'),
  },
  {
    field: 'created-at',
    headerName: 'Created At',
    resizable: false,
    flex: 1,
    minWidth: 150,
    renderCell: (params: GridRenderCellParams) => {
      return format(new Date(params.row['created-at']), 'dd.MM.yyyy HH:mm:ss');
    },
  },
  {
    field: 'updated-at',
    headerName: 'Updated At',
    resizable: false,
    flex: 1,
    minWidth: 150,
    renderCell: (params: GridRenderCellParams) => {
      return format(new Date(params.row['updated-at']), 'dd.MM.yyyy HH:mm:ss');
    },
  },
  {
    field: 'deleted-at',
    headerName: 'Deleted At',
    resizable: false,
    flex: 1,
    minWidth: 150,
    renderCell: (params: GridRenderCellParams) => {
      return params.row['deleted-at'] 
        ? format(new Date(params.row['deleted-at']), 'dd.MM.yyyy HH:mm:ss')
        : '-';
    },
  },
  {
    field: 'subscription-start',
    headerName: 'Subscription Start',
    resizable: false,
    flex: 1,
    minWidth: 150,
    renderCell: (params: GridRenderCellParams) => {
      return format(new Date(params.row['subscription-start']), 'dd.MM.yyyy HH:mm:ss');
    },
  },
  {
    field: 'subscription-end',
    headerName: 'Subscription End',
    resizable: false,
    flex: 1,
    minWidth: 150,
    renderCell: (params: GridRenderCellParams) => {
      return format(new Date(params.row['subscription-end']), 'dd.MM.yyyy HH:mm:ss');
    },
  },
  {
    field: 'actions',
    type: 'actions',
    headerName: 'Actions',
    resizable: false,
    flex: 1,
    minWidth: 80,
    getActions: () => {
      return [
        <Tooltip title="Edit">
          <GridActionsCellItem
            icon={
              <IconifyIcon
                icon="fluent:edit-32-filled"
                color="text.secondary"
                sx={{ fontSize: 'body1.fontSize', pointerEvents: 'none' }}
              />
            }
            label="Edit"
            size="small"
          />
        </Tooltip>,
        <Tooltip title="Delete">
          <GridActionsCellItem
            icon={
              <IconifyIcon
                icon="mingcute:delete-3-fill"
                color="error.main"
                sx={{ fontSize: 'body1.fontSize', pointerEvents: 'none' }}
              />
            }
            label="Delete"
            size="small"
          />
        </Tooltip>,
      ];
    },
  },
];

const CustomerTable = ({ searchText }: { searchText: string }): ReactElement => {
  const apiRef = useGridApiRef<GridApi>();

  const visibleColumns = useMemo(
    () => columns.filter((column) => column.field !== 'id'),
    [columns],
  );

  useEffect(() => {
    apiRef.current.setQuickFilterValues(
      searchText.split(/\b\W+\b/).filter((word: string) => word !== ''),
    );
  }, [searchText]);

  useEffect(() => {
    const handleResize = () => {
      if (apiRef.current) {
        apiRef.current.resize();
      }
    };
    window.addEventListener('resize', handleResize);
    return () => {
      window.removeEventListener('resize', handleResize);
    };
  }, [apiRef]);

  return (
    <>
      <DataGrid
        apiRef={apiRef}
        density="standard"
        columns={visibleColumns}
        autoHeight={false}
        rowHeight={56}
        checkboxSelection
        disableColumnMenu
        disableRowSelectionOnClick
        rows={rows}
        onResize={() => {
          apiRef.current.autosizeColumns({
            includeOutliers: true,
            expand: true,
          });
        }}
        initialState={{
          pagination: { paginationModel: { page: 0, pageSize: 4 } },
        }}
        slots={{
          loadingOverlay: CircularProgress as GridSlots['loadingOverlay'],
          pagination: CustomPagination as GridSlots['pagination'],
          noResultsOverlay: CustomNoResultsOverlay as GridSlots['noResultsOverlay'],
        }}
        slotProps={{
          pagination: { labelRowsPerPage: rows.length },
        }}
        sx={{
          height: 1,
          width: 1,
          tableLayout: 'fixed',
          scrollbarWidth: 'thin',
        }}
      />
    </>
  );
};

export default CustomerTable;
