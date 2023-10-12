import React from "react";
import { User, UserPermission } from "./generated";
import { API } from "./state";
import { DataGrid } from "@mui/x-data-grid";
import { Box } from "@mui/material";
import { useSnackbar } from "notistack";

export const Users = () => {
  const [users, setUsers] = React.useState<User[]>([]);
  React.useEffect(() => {
    API.getUsers().then(setUsers);
  }, []);
  const { enqueueSnackbar } = useSnackbar();

  return (
    <Box display="flex" justifyContent="center">
      <Box maxWidth={1000} width="100%" height={500}>
        <DataGrid<User>
          rows={users}
          columns={[
            { field: "id", headerName: "ID", width: 100 },
            { field: "username", headerName: "Name", width: 300 },
            {
              field: "permission",
              headerName: "Permission",
              width: 150,
              type: "singleSelect",
              valueOptions: Object.values(UserPermission),
              valueSetter: (p) => {
                API.setUserPermission({
                  permission: p.value as any,
                  userID: p.row.id,
                })
                  .then(() => {
                    enqueueSnackbar("permission changed", {
                      variant: "success",
                    });
                  })
                  .catch(() => {
                    enqueueSnackbar("something went wrong", {
                      variant: "error",
                    });
                  });
                return { ...p.row, permission: p.value };
              },

              editable: true,
            },
            {
              field: "createdAt",
              headerName: "Created At",
              valueGetter: (x) => x.value.toISOString(),
              width: 200,
            },
          ]}
        ></DataGrid>
      </Box>
    </Box>
  );
};
