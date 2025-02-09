import React from 'react';
import {
  Card,
  CardContent,
  Typography,
  IconButton,
} from '@mui/material';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import { green, red, grey } from '@mui/material/colors';
import { ContainerInfo } from '../model/ContainerInfo';
import { ContainerInfoGrid } from './ContainerInfoGrid';

interface ContainerCardProps {
  container: ContainerInfo;
}

export const ContainerCard = ({ container } : ContainerCardProps) => (
    <Card variant="outlined" sx={{ margin: 2 }}>
        <CardContent>
            <div style={{
                display: "flex",
                gap: "10px",
                marginLeft: "20px",
                alignItems: "center",
                fontSize: "15px"
            }}>
                <div style={{
                width: "20px",
                height: "20px",
                backgroundColor: container.status === 'online' ? green[400] : red[400],
                borderRadius: "25px"
                }} />
                <Typography variant="h6" component="h3">{container.name}</Typography>
            </div>

            <div style={{
                display: "flex",
                gap: "5px",
                marginLeft: "15px",
                marginTop: "5px",
                alignItems: "center"
            }}>
                <IconButton  size="small" onClick={() => navigator.clipboard.writeText(container.id)}>
                    <ContentCopyIcon fontSize="small" />
                </IconButton>
                
                <Typography variant="subtitle2" color={grey[700]}>
                    {container.dockerId}
                </Typography>
            </div>

            <ContainerInfoGrid container={container} />
        </CardContent>
    </Card>
    );