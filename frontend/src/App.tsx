import { useEffect, useState } from 'react';
import './App.css';
import { ContainerCard } from './components/ContainerCard';
import { ContainerInfo } from './model/ContainerInfo';
import { REFRESH_RATE_SEC } from './values';
import axios from 'axios';
import { Typography } from '@mui/material';

function App() {
    const [containers, setContainers] = useState<ContainerInfo[]>([]);

    const fetchContainers = async () => {
        axios.get<ContainerInfo[]>("/container", { params: { offset: 0, limit: 1000 } })
            .then(resp => setContainers(resp.data))
            .catch(err => console.log(err));
    }

    useEffect(() => {
        let id = setInterval(fetchContainers, REFRESH_RATE_SEC * 1000);
        return () => clearInterval(id);
    }, [])

    return (
        <div style={{maxWidth: "1100px", marginLeft: 'auto', marginRight: 'auto', marginTop: '50px'}}>
            
            {containers.map((container: ContainerInfo, index: number) => (
                <div key={container.id} style={{ padding: '5px'}}>
                <ContainerCard container={container}/>
            </div>
            ))}

            {containers.length === 0 && 
                <Typography variant='h4' align='center'>No containers found</Typography>
            }
        </div>
    );
}

export default App;
