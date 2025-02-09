import { Grid, Paper, Typography, Chip } from '@mui/material';
import { ContainerInfo } from '../model/ContainerInfo';
import { grey } from '@mui/material/colors';

interface ContainerInfoGridProps {
    container: ContainerInfo;
}

export const ContainerInfoGrid = ({ container } : ContainerInfoGridProps) => (
  <Paper elevation={0} sx={{ p: 2 }}>
    <Grid container spacing={3}>
      <Grid item xs={2}>
        <Typography variant="subtitle2" color={grey[700]}>Name</Typography>
        <Typography>{container.name}</Typography>
      </Grid>
      <Grid item xs={2}>
        <Typography variant="subtitle2" color={grey[700]}>IP</Typography>
        <Typography>{container.ip}</Typography>
      </Grid>
      <Grid item xs={3}>
        <Typography variant="subtitle2" color={grey[700]}>Last Check</Typography>
        <Typography>{container.lastCheck}</Typography>
      </Grid>
      <Grid item xs={3}>
        <Typography variant="subtitle2" color={grey[700]}>Latest Activity</Typography>
        <Typography>{container.lastActivity}</Typography>
      </Grid>
      <Grid item xs={1}>
        <Typography variant="subtitle2" color={grey[700]}>Status</Typography>
        <Chip 
          label={container.status} 
          color={container.status === 'online' ? 'success' : 'error'} 
          size="small"
        />
      </Grid>
    </Grid>
  </Paper>
);