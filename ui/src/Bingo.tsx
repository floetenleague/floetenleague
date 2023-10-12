import { Box, Button, Tooltip, Typography, useTheme } from "@mui/material";
import { chunk } from "lodash";
import {
  FLBingoField,
  FLBingoFieldStatus,
  FLBingoUserBoard,
  UserPermission,
} from "./generated";
import { useParams } from "react-router-dom";
import { useStore } from "./state";
import { useSnackbar } from "notistack";

export interface MiniBoardProps {
  fields: FLBingoField[];
  userFields: FLBingoUserBoard["fields"];
  columns: number;
}

interface StatusStuffItem {
  background: string;
}

const statusStuff: Record<FLBingoFieldStatus, StatusStuffItem> = {
  [FLBingoFieldStatus.DoneInReview]: {
    background: "#e67e22",
  },
  [FLBingoFieldStatus.Done]: {
    background: "#27ae60",
  },
  [FLBingoFieldStatus.Blank]: {
    background: "#3c3836",
  },
  [FLBingoFieldStatus.Bingo]: {
    background: "#f1c40f",
  },
};

export interface NamedMiniBoardProps extends MiniBoardProps {
  name: string;
  score: number;
}
export const NamedMiniBoard = ({
  name,
  score,
  ...miniBoard
}: NamedMiniBoardProps) => {
  const theme = useTheme();
  return (
    <Box
      style={{ background: theme.palette.background.paper, cursor: "pointer" }}
      padding={0.5}
      display="inline-block"
    >
      <Typography variant="caption" component="div" textAlign="center">
        <b>[{score}P]</b> {name}
      </Typography>
      <MiniBoard {...miniBoard} />
    </Box>
  );
};

export const MiniBoard = ({ columns, fields, userFields }: MiniBoardProps) => {
  const gridGap = 5;

  return (
    <Box
      display="grid"
      gridAutoColumns="1fr"
      width={200}
      height={200}
      style={{ gridGap, padding: 2 }}
    >
      {chunk(fields, columns).map((chunk, index) => (
        <Box
          key={index}
          style={{ gridGap }}
          display="grid"
          gridAutoFlow="column"
        >
          {chunk.map((field) => (
            <Box
              data-id={field.id}
              key={field.id}
              style={{
                background:
                  statusStuff[
                    userFields[field.id]?.status ?? FLBingoFieldStatus.Blank
                  ].background,
                maxHeight: "100%",
              }}
            />
          ))}
        </Box>
      ))}
    </Box>
  );
};

export interface BoardProps {
  name: string;
  fields: FLBingoField[];
  userFields: FLBingoUserBoard["fields"];
  columns: number;
}

export const BoardPage = () => {
  const { id, userId } = useParams();
  const { enqueueSnackbar } = useSnackbar();
  const bingo = useStore((state) => state.overview).bingos.find(
    (bingo) => bingo.id === parseInt(id ?? "0")
  );
  const board = bingo?.boards.find(
    (board) => board.userId === parseInt(userId ?? "0")
  );
  const auth = useStore((state) => state.auth);
  const _setFieldStatus = useStore((state) => state.setFieldStatus);
  const setFieldStatus = (
    fieldID: number,
    status: FLBingoFieldStatus,
    success: string
  ) => {
    _setFieldStatus({
      bingoID: bingo?.id!!,
      userID: parseInt(userId ?? "0"),
      status,
      fieldID,
    })
      .then(() => {
        enqueueSnackbar(success, { variant: "success" });
      })
      .catch(() => {
        enqueueSnackbar("could not update field status :/", {
          variant: "error",
        });
      });
  };

  const theme = useTheme();
  const gridGap = 10;
  if (bingo == null || board == null) {
    return <>not found</>;
  }

  return (
    <Box display="flex" justifyContent="center">
      <Box
        style={{ background: theme.palette.background.paper }}
        padding={1}
        display="inline-block"
        maxWidth={1000}
        minWidth={700}
        width="100%"
      >
        <Typography variant="h3" component="div" textAlign="center">
          [{board.score}P] {board.username}
        </Typography>
        <Box
          display="grid"
          maxWidth={5000}
          gridAutoRows={`minmax(0, 1fr)`}
          width="100%"
          style={{ gridGap, padding: 2, aspectRatio: "1/1" }}
        >
          {chunk(bingo.fields, bingo.size).map((chunk, index) => (
            <Box
              key={index}
              style={{ gridGap }}
              display="grid"
              gridAutoColumns={`minmax(0, 1fr)`}
              gridAutoFlow="column"
            >
              {chunk.map((field) => {
                const status =
                  board.fields[field.id]?.status ?? FLBingoFieldStatus.Blank;
                const { background } = statusStuff[status];
                const color = theme.palette.getContrastText(background);
                return (
                  <Tooltip
                    arrow
                    title={
                      <>
                        <Typography marginBottom={1} fontWeight="bold">
                          {field?.label}
                        </Typography>
                        <Typography>{field.description}</Typography>
                      </>
                    }
                  >
                    <Box
                      key={field.id}
                      display="flex"
                      flexDirection="column"
                      color={color}
                      style={{ background }}
                    >
                      <Typography
                        padding={1}
                        flexGrow="1"
                        height={40}
                        textOverflow="ellipsis"
                        overflow="hidden"
                      >
                        <div
                          style={{
                            overflow: "hidden",
                            textOverflow: "ellipsis",
                          }}
                        >
                          [{field.score}P] {field?.label}
                        </div>
                      </Typography>
                      <Box padding={1}>
                        {auth.permission === UserPermission.Moderator ? (
                          <BoardModAction
                            fieldID={field.id}
                            status={status}
                            setFieldStatus={setFieldStatus}
                          />
                        ) : auth.id === board.userId ? (
                          <BoardUserAction
                            fieldID={field.id}
                            status={status}
                            setFieldStatus={setFieldStatus}
                          />
                        ) : undefined}
                      </Box>
                    </Box>
                  </Tooltip>
                );
              })}
            </Box>
          ))}
        </Box>
      </Box>
    </Box>
  );
};

const BoardModAction = ({
  status,
  fieldID,
  setFieldStatus,
}: {
  fieldID: number;
  status: FLBingoFieldStatus;
  setFieldStatus: (
    fieldID: number,
    status: FLBingoFieldStatus,
    success: string
  ) => void;
}) => {
  switch (status) {
    case FLBingoFieldStatus.Blank:
      return (
        <Button
          fullWidth
          variant="contained"
          onClick={() => {
            setFieldStatus(fieldID, FLBingoFieldStatus.Done, "Field updated");
          }}
          size="small"
        >
          done
        </Button>
      );
    case FLBingoFieldStatus.DoneInReview:
      return (
        <>
          <Button
            fullWidth
            variant="contained"
            onClick={() => {
              setFieldStatus(
                fieldID,
                FLBingoFieldStatus.Blank,
                "Field updated"
              );
            }}
            size="small"
          >
            Set blank
          </Button>
          <Button
            fullWidth
            variant="contained"
            onClick={() => {
              setFieldStatus(fieldID, FLBingoFieldStatus.Done, "Field updated");
            }}
            size="small"
          >
            done
          </Button>
        </>
      );
    case FLBingoFieldStatus.Done:
    case FLBingoFieldStatus.Bingo:
      return (
        <Button
          fullWidth
          variant="contained"
          onClick={() => {
            setFieldStatus(fieldID, FLBingoFieldStatus.Blank, "Field updated");
          }}
          size="small"
        >
          undo
        </Button>
      );
    default:
      return <></>;
  }
};
const BoardUserAction = ({
  status,
  fieldID,
  setFieldStatus,
}: {
  fieldID: number;
  status: FLBingoFieldStatus;
  setFieldStatus: (
    fieldID: number,
    status: FLBingoFieldStatus,
    success: string
  ) => void;
}) => {
  switch (status) {
    case FLBingoFieldStatus.Blank:
      return (
        <Button
          fullWidth
          variant="contained"
          onClick={() => {
            setFieldStatus(
              fieldID,
              FLBingoFieldStatus.DoneInReview,
              "Updated successfully. A moderator will review your progress."
            );
          }}
          size="small"
        >
          done
        </Button>
      );
    case FLBingoFieldStatus.DoneInReview:
      return (
        <>
          <Typography textAlign="center">In Mod Review</Typography>
          <Button
            fullWidth
            variant="contained"
            size="small"
            onClick={() => {
              setFieldStatus(
                fieldID,
                FLBingoFieldStatus.Blank,
                "Updated successfully"
              );
            }}
          >
            Cancel
          </Button>
        </>
      );
    default:
      return <></>;
  }
};
