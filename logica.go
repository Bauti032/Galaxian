package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	constCantFilasTablero    = 30
	constCantColumnasTablero = 30

	constCantColumnas = 2
	constY            = 0
	constX            = 1

	constCantColumnasOvni = 4
	constTipoOvni         = 0
	constOvniY            = 1
	constOvniX            = 2
	constEnDescenso       = 3

	constTiempoDeDisparoOvni   = 3
	constTiempoLiberarcionOvni = 10

	constSimboloVacío       = ""
	constSimboloNave        = "N"
	constSimboloDisparoNave = "*"
	constSimboloDisparoOvni = "."
	constSimboloOvniLider   = "L"
	constSimboloOvniComun   = "C"
	constSimboloBorde       = "X"

	constCantColumnasDisparos = 2
)

// Vector global con las direcciones posibles
var (
	quieto    = [constCantColumnas]int{0, 0}
	izquierda = [constCantColumnas]int{0, -1}
	derecha   = [constCantColumnas]int{0, 1}
	arriba    = [constCantColumnas]int{-1, 0}
	abajo     = [constCantColumnas]int{1, 0}
)

// Vector global con las direccion de la nave
var direccionNave [constCantColumnas]int

// Variable global que indica si se presiono la barra espaciadora lo que ejecuta un disparo de la nave
var disparoNave bool

// Función para enviar actualizaciones a los clientes
func generarEventos() {
	var (
		tablero [constCantFilasTablero][constCantColumnasTablero]string

		nave         [constCantColumnas]int
		disparosNave [][constCantColumnas]int

		ovnis         [][constCantColumnasOvni]int
		disparosOvnis [][constCantColumnas]int

		ultimaEjecucionDisparoOvni    time.Time
		ultimaEjecucionLiberacionOvni time.Time

		puntos int
		vidas  int
	)

	rand.Seed(time.Now().Unix())

	//Se inicializa variables
	ultimaEjecucionDisparoOvni = time.Now()
	ultimaEjecucionLiberacionOvni = time.Now()

	disparoNave = false

	vidas = 3

	// Se genera tablero por primera vez con los bordes
	tablero = generarTablero()

	// Se genera la nave (posición inicial) por primera vez
	nave, direccionNave = inicializarNave(constCantFilasTablero, constCantColumnasTablero)

	// Se generan los ovnis (posiciones iniciales) por primera vez
	ovnis = inicializarOvnis(constCantFilasTablero, constCantColumnasTablero)

	// Se actualiza nave y ovnis en el tablero por primera vez
	actualizarTablero(&tablero, nave, disparosNave, ovnis, disparosOvnis)

	for {
		// Se actualizan las posiciones de la nave según la dirección
		calcularNuevaPosicionNave(tablero, &nave, &direccionNave)

		// Se crea un nuevo disparo si corresponde
		crearDisparoNave(nave, &disparoNave, &disparosNave)

		//Cada "constTiempoDeDisparoOvni" segundos, se crea un disparo de un ovni
		if time.Since(ultimaEjecucionDisparoOvni) >= constTiempoDeDisparoOvni*time.Second {
			crearDisparoOvni(ovnis, &disparosOvnis)
			ultimaEjecucionDisparoOvni = time.Now()
		}

		//Cada "constTiempoLiberarcionOvni" segundos, se libera un obvni de la formación
		if time.Since(ultimaEjecucionLiberacionOvni) >= constTiempoLiberarcionOvni*time.Second {
			liberarOvni(ovnis)
			ultimaEjecucionLiberacionOvni = time.Now()
		}

		// Se calcula la nueva posición de los ovnis liberados
		calcularNuevaPosicionOvnisLiberados(ovnis)

		// Se calcula las nuevas posiciones de los disparos de la nave y de los ovnis
		calcularNuevasPosicionesDisparos(tablero, disparosNave, disparosOvnis)

		// Se verifica el estado del juego y eliminan elementos si corresponde
		if !verificarEstadoDeJuego(tablero, nave, &ovnis, &disparosNave, &disparosOvnis, &puntos) {
			// Si no tiene más vidas, se devuelve pantalla gameOver
			vidas--

			if vidas <= 0 {
				enviarGameOver(puntos)
				return
			}
		} else {
			if len(ovnis) == 0 {
				enviarWin(puntos)

				return
			}

			enviarActualizacionTexto(fmt.Sprint("Puntaje: ", puntos, ". Vidas: ", vidas))
		}

		//Se actualiza el tablero con los valores de la nave, ovnis y disparos en sus nuevas posiciones
		actualizarTablero(&tablero, nave, disparosNave, ovnis, disparosOvnis)

		// Se envía actualización de tablero al cliente para mostrar en pantalla
		enviarActualizacionTablero(tablero)

		// Espera un tiempo antes de generar un nuevo movimiento
		time.Sleep(85 * time.Millisecond)
	}
}

func generarTablero() [constCantFilasTablero][constCantColumnasTablero]string {
	var tablero [constCantFilasTablero][constCantColumnasTablero]string

	for x := 0; x < constCantColumnasTablero; x++ {
		if x == 0 || x == 29 {
			for y := 0; y < constCantFilasTablero; y++ {
				tablero[y][x] = constSimboloBorde
			}
		}
	}

	for y := 0; y < constCantFilasTablero; y++ {
		if y == 0 || y == 29 {
			for x := 0; x < constCantColumnasTablero; x++ {
				tablero[y][x] = constSimboloBorde
			}
		}
	}

	return tablero
}

func inicializarNave(cantFilasTablero int, cantColumnasTablero int) ([constCantColumnas]int, [constCantColumnas]int) {

	var (
		nave [constCantColumnas]int
	)

	nave[constY] = cantFilasTablero - 4
	nave[constX] = cantColumnasTablero / 2

	return nave, quieto
}

func inicializarOvnis(cantFilasTablero int, cantColumnasTablero int) [][constCantColumnasOvni]int {
	var (
		vec   [4]int
		ovnis [][constCantColumnasOvni]int
	)
	rand.Seed(time.Now().Unix())

	for y := 2; y < cantFilasTablero-20; y++ {
		for x := 5; x < cantColumnasTablero-5; x++ {

			tipo := rand.Intn(2) + 1

			vec[0] = tipo
			vec[1] = y
			vec[2] = x
			vec[3] = 0

			ovnis = append(ovnis, vec)

		}
	}

	return ovnis
}

func actualizarTablero(tablero *[constCantFilasTablero][constCantColumnasTablero]string,
	nave [constCantColumnas]int,
	disparosNave [][constCantColumnas]int,
	ovnis [][constCantColumnasOvni]int,
	disparosOvnis [][constCantColumnas]int) {

	//bucle for para limpiar el tablero
	for y := 1; y < constCantFilasTablero-1; y++ {
		for x := 1; x < constCantColumnasTablero-1; x++ {
			tablero[y][x] = constSimboloVacío
		}
	}

	//inicializar nave
	tablero[nave[constY]][nave[constX]] = constSimboloNave

	//inicio ovnis
	for _, ovni := range ovnis {
		var simbolo string
		if ovni[constTipoOvni] == 1 {
			simbolo = constSimboloOvniLider
		} else {
			simbolo = constSimboloOvniComun
		}
		tablero[ovni[constOvniY]][ovni[constOvniX]] = simbolo
	}

	for y := 0; y < len(disparosNave); y++ {
		tablero[disparosNave[y][constY]][disparosNave[y][constX]] = constSimboloDisparoNave
	}

	for y := 0; y < len(disparosOvnis); y++ {
		tablero[disparosOvnis[y][constY]][disparosOvnis[y][constX]] = constSimboloDisparoOvni
	}

}

func calcularNuevaPosicionNave(tablero [constCantFilasTablero][constCantColumnasTablero]string,
	nave *[constCantColumnas]int, direccionNave *[constCantColumnas]int) {

	nuevaY := nave[constY] + direccionNave[constY]
	nuevaX := nave[constX] + direccionNave[constX]

	// Verifica que la nueva posición no sea un borde
	if tablero[nuevaY][nuevaX] != constSimboloBorde {
		nave[constY] = nuevaY
		nave[constX] = nuevaX
	}
}

func crearDisparoNave(nave [constCantColumnas]int,
	disparoNave *bool,
	disparosNave *[][constCantColumnasDisparos]int) {

	if *disparoNave {
		disparo := [constCantColumnasDisparos]int{
			nave[constY] - 1, // Justo arriba de la nave
			nave[constX],
		}
		*disparosNave = append(*disparosNave, disparo)
		*disparoNave = false // Para evitar disparos continuos
	}

}

func crearDisparoOvni(ovnis [][constCantColumnasOvni]int,
	disparosOvnis *[][constCantColumnasDisparos]int) {

	indice := rand.Intn(len(ovnis))
	ovni := ovnis[indice]

	disparo := [constCantColumnasDisparos]int{
		ovni[constOvniY] + 1,
		ovni[constOvniX],
	}
	*disparosOvnis = append(*disparosOvnis, disparo)
}

func calcularNuevasPosicionesDisparos(tablero [constCantFilasTablero][constCantColumnasTablero]string,
	disparosNave [][constCantColumnasDisparos]int,
	disparosOvnis [][constCantColumnasDisparos]int) {

	for y := 0; y < len(disparosNave); y++ {
		disparosNave[y][0] -= 1
	}

	for y := 0; y < len(disparosOvnis); y++ {
		disparosOvnis[y][0] += 1
	}

}

func verificarEstadoDeJuego(
	tablero [constCantFilasTablero][constCantColumnasTablero]string,
	nave [constCantColumnas]int,
	ovnis *[][constCantColumnasOvni]int,
	disparosNave *[][constCantColumnasDisparos]int,
	disparosOvnis *[][constCantColumnasDisparos]int,
	puntos *int,
) bool {

	// 1. Disparos de la nave
	for i := 0; i < len(*disparosNave); {
		y := (*disparosNave)[i][constY]
		x := (*disparosNave)[i][constX]

		switch tablero[y][x] {
		case constSimboloBorde:
			// Toca borde, eliminar disparo
			*disparosNave = eliminarDisparo(*disparosNave, y, x)
		case constSimboloDisparoOvni:
			// Toca disparo ovni, eliminar disparo nave
			*disparosNave = eliminarDisparo(*disparosNave, y, x)
		case constSimboloOvniComun:
			// Toca ovni común no en descenso
			for _, ovni := range *ovnis {
				if ovni[constOvniY] == y && ovni[constOvniX] == x && ovni[constTipoOvni] == 2 && ovni[constEnDescenso] == 0 {
					*disparosNave = eliminarDisparo(*disparosNave, y, x)
					*ovnis = eliminarOvni(*ovnis, y, x)
					*puntos += 10
					break
				}
			}
		case constSimboloOvniLider:
			// Toca ovni líder no en descenso
			for j, ovni := range *ovnis {
				if ovni[constOvniY] == y && ovni[constOvniX] == x && ovni[constTipoOvni] == 1 && ovni[constEnDescenso] == 0 {
					*disparosNave = eliminarDisparo(*disparosNave, y, x)
					(*ovnis)[j][constTipoOvni] = 2 // Convertir a común
					*puntos += 20
					break
				}
			}
		default:
			i++
		}
	}

	// 2. Disparos de ovni
	for i := 0; i < len(*disparosOvnis); {
		y := (*disparosOvnis)[i][constY]
		x := (*disparosOvnis)[i][constX]

		switch tablero[y][x] {
		case constSimboloBorde:
			// Toca borde, eliminar disparo ovni
			*disparosOvnis = eliminarDisparo(*disparosOvnis, y, x)
		case constSimboloDisparoNave:
			// Toca disparo nave, eliminar disparo nave
			*disparosNave = eliminarDisparo(*disparosNave, y, x)
			i++
		case constSimboloNave:
			// Toca nave, pierde vida
			return false
		default:
			i++
		}
	}

	// 3. Ovni choca nave o borde
	for i := 0; i < len(*ovnis); {
		y := (*ovnis)[i][constOvniY]
		x := (*ovnis)[i][constOvniX]

		if tablero[y][x] == constSimboloNave {
			*ovnis = eliminarOvni(*ovnis, y, x)
			return false
		} else if tablero[y][x] == constSimboloBorde {
			*ovnis = eliminarOvni(*ovnis, y, x)
		} else {
			i++
		}
	}

	return true
}

func eliminarDisparo(slice [][constCantColumnasDisparos]int, coordenadaY int, coordenadaX int) [][2]int {
	var nuevoSlice [][constCantColumnasDisparos]int
	for f := 0; f < len(slice); f++ {
		if slice[f][constY] != coordenadaY &&
			slice[f][constX] != coordenadaX {
			nuevoSlice = append(nuevoSlice, slice[f])
		}
	}
	return nuevoSlice
}

func eliminarOvni(slice [][constCantColumnasOvni]int, coordenadaY int, coordenadaX int) [][4]int {
	var nuevoSlice [][constCantColumnasOvni]int
	for f := 0; f < len(slice); f++ {
		if slice[f][constOvniY] != coordenadaY ||
			slice[f][constOvniX] != coordenadaX {
			nuevoSlice = append(nuevoSlice, slice[f])
		}
	}
	return nuevoSlice
}

func liberarOvni(ovnis [][constCantColumnasOvni]int) {
	indice := rand.Intn(len(ovnis))
	ovnis[indice][constEnDescenso] = 1

}

func calcularNuevaPosicionOvnisLiberados(ovnis [][constCantColumnasOvni]int) {

	for y := 0; y < len(ovnis); y++ {
		if ovnis[y][3] == 1 {
			ovnis[y][1] += 1

			if ovnis[y][1] == 29 {

				ovnis = eliminarOvni(ovnis, ovnis[y][1], ovnis[y][2])
			}

		}
	}

}
