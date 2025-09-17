-- Script simplificado de inicialización para PostgreSQL en Docker

-- Tabla principal de pilotos
CREATE TABLE IF NOT EXISTS pilotos (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL UNIQUE,
    equipo VARCHAR(255),
    nacionalidad VARCHAR(100),
    numero INT,
    victorias INT,
    puntos NUMERIC,
    podios INT,
    poles INT,
    vueltas_rapidas INT
);

-- Insertar pilotos de ejemplo
INSERT INTO pilotos (nombre, equipo, nacionalidad, numero, victorias, puntos, podios, poles, vueltas_rapidas) VALUES
('Max Verstappen', 'Red Bull', 'Holandés', 1, 54, 575.5, 98, 33, 30),
('Lewis Hamilton', 'Mercedes', 'Británico', 44, 103, 4637.5, 197, 104, 65),
('Charles Leclerc', 'Ferrari', 'Monegasco', 16, 5, 1074, 32, 23, 9),
('Fernando Alonso', 'Aston Martin', 'Español', 14, 32, 2267, 106, 22, 24),
('Carlos Sainz', 'Ferrari', 'Español', 55, 2, 782.5, 18, 5, 3)
ON CONFLICT (nombre) DO NOTHING;

-- Crear índices básicos
CREATE INDEX IF NOT EXISTS idx_pilotos_equipo ON pilotos(equipo);
CREATE INDEX IF NOT EXISTS idx_pilotos_nacionalidad ON pilotos(nacionalidad);
