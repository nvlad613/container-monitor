
export interface ContainerInfo {
    name: string;
    ip: string;
    id: string;
    dockerId: string,
    lastCheck: string;
    lastActivity: string;
    status: 'online' | 'offline'; // Статус может быть только 'online' или 'offline'
}