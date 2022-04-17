package servidor

import (
	"crud/banco"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type usuario struct {
	ID    uint32 `json:"id"`
	NOME  string `json:"nome"`
	EMAIL string `json:"email"`
}

func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	corpoRequisicao, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.Write([]byte("Falha ao ler o corpo da requisição!"))
		return
	}

	var usuario usuario

	if err = json.Unmarshal(corpoRequisicao, &usuario); err != nil {
		w.Write([]byte("Erro ao converter o usuário para struct"))
		return
	}

	db, err := banco.Conectar()

	if err != nil {
		w.Write([]byte("Erro ao conectar no banco de dados"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("insert into usuarios (nome, email) values ($1, $2)")

	if err != nil {
		w.Write([]byte("Erro ao criar statement"))
		return
	}

	defer statement.Close()

	insercao, err := statement.Exec(usuario.NOME, usuario.EMAIL)

	if err != nil {
		w.Write([]byte("Erro ao executar o statement"))
		return
	}

	idInserido, err := insercao.LastInsertId()

	if err != nil {
		w.Write([]byte("Erro ao obter id inserido!"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Usuário inserido com sucesso! Id: %d", idInserido)))

	// fmt.Println(usuario)
}

func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	db, err := banco.Conectar()

	if err != nil {
		w.Write([]byte("Erro ao conectar com banco de dados"))
	}
	defer db.Close()

	linhas, err := db.Query("select * from usuarios")

	if err != nil {
		w.Write([]byte("Erro ao buscar usuários"))
		return
	}
	defer linhas.Close()

	var usuarios []usuario
	for linhas.Next() {
		var usuario usuario

		if err := linhas.Scan(&usuario.ID, &usuario.NOME, &usuario.EMAIL); err != nil {
			w.Write([]byte("Erro ao escanear o usuário"))
			return
		}

		usuarios = append(usuarios, usuario)
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(usuarios); err != nil {
		w.Write([]byte("Erro ao converter os usuários para JSON"))
		return
	}
}

func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	ID, err := strconv.ParseUint(parametros["id"], 10, 32)

	if err != nil {
		w.Write([]byte("Erro ao converter o parâmentro para inteiro"))
		return
	}

	db, err := banco.Conectar()

	if err != nil {
		w.Write([]byte("Erro ao conectar"))
		return
	}

	linha, err := db.Query("select * from usuarios where id = $1", ID)

	if err != nil {
		w.Write([]byte("Erro ao buscar usuário"))
		return
	}

	var usuario usuario

	if linha.Next() {
		if err := linha.Scan(&usuario.ID, &usuario.NOME, &usuario.EMAIL); err != nil {
			w.Write([]byte("Erro ao escanear usuário"))
			return
		}
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(usuario); err != nil {
		w.Write([]byte("Erro ao converter usuário para JSON"))
		return
	}
}

func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	ID, err := strconv.ParseUint(parametros["id"], 10, 32)

	if err != nil {
		w.Write([]byte("Erro ao converter o parâmetro para inteiro"))
		return
	}

	corpoRequisicao, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.Write([]byte("Erro ao ler o corpo da requisição"))
		return
	}

	var usuario usuario

	if err := json.Unmarshal(corpoRequisicao, &usuario); err != nil {
		w.Write([]byte("Erro ao converter usuário para struct"))
		return
	}

	db, err := banco.Conectar()

	if err != nil {
		w.Write([]byte("Erro ao conectar com banco de dados"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("update usuarios set nome = $1, email = $2 where id = $3 ")

	if err != nil {
		w.Write([]byte("Erro ao criar statement"))
		return
	}
	defer statement.Close()

	if _, err := statement.Exec(&usuario.NOME, &usuario.EMAIL, ID); err != nil {
		w.Write([]byte("Erro ao atulizar o usuário"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	ID, err := strconv.ParseUint(parametros["id"], 10, 32)

	if err != nil {
		w.Write([]byte("Erro ao converter o parâmetro para inteiro"))
		return
	}

	db, err := banco.Conectar()

	if err != nil {
		w.Write([]byte("Erro ao conectar"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("delete from usuarios where id = $1")

	if err != nil {
		w.Write([]byte("Erro ao criar statement"))
		return
	}
	defer statement.Close()

	if _, err := statement.Exec(ID); err != nil {
		w.Write([]byte("Erro ao deletar usuário"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
