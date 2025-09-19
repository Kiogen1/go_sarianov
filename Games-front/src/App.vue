<template>
  <div id="app" class="container">
    <h1>Игры</h1>

    <!-- форма добавления -->
    <div class="form">
      <input v-model="newGame.name" placeholder="Название" />
      <input v-model="newGame.studio" placeholder="Студия разработчик" />
      <input v-model.number="newGame.year" placeholder="Год выпуска" />
      <input v-model.number="newGame.sold" placeholder="Продано копий" />
      <button @click="saveGame">Добавить</button>
    </div>
    <div v-if="errorMessage.length">
      <ul>
        <li v-for="err in errorMessage" :key="err">{{ err }}</li>
      </ul>
    </div>

    <!-- таблица -->
    <table border="1" cellpadding="5">
      <thead>
        <tr>
          <th>ID</th>
          <th>Название</th>
          <th>Студия</th>
          <th>Год основания</th>
          <th>Продано копий</th>
          <th>Кнопки</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="game in strGames" :key="game.id">
          <td>{{ game.id }}</td>
          <td>{{ game.name }}</td>
          <td>{{ game.studio }}</td>
          <td>{{ game.year }}</td>
          <td>{{ game.sold }}</td>
          <td>
            <button @click="editGame(game)">Редактировать</button>
            <button @click="deleteGame(game.id)">Удалить</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
//import axios from "axios";
import api from "./api.js"

export default {
  name: "App",
  data() {
    return {
      strGames: [],
      newGame: { name: "", studio: "", year: "", sold: "" },
      editing: null,
      errorMessage: [],
      count: 0
    };
  },
  created() {
    this.fetchGames();
  },
  methods: {
    validateGame(game) {
      const errors = [];
      // Пустые поля
      if(!game.name) errors.push("Поле 'Название' не должно быть пустым");
      if(!game.studio) errors.push("Поле 'Страна' не должно быть пустым");
      if(!game.year) errors.push("Поле 'Год' не должно быть пустым");
      if(!game.sold) errors.push("Поле 'Капитализация' не должно быть пустым")

      // Проверка чисел
      const currentDate = new Date();
      const currentYear = currentDate.getFullYear();
      if(game.year && isNaN(Number(game.year)) || Number(game.year) > currentYear) errors.push("Поле 'Год' может состоять только из цифр, и не быть больше текущего года");
      if(game.sold && isNaN(Number(game.sold))) errors.push("Поле 'Продано копий' может состоять только из цифр");

      // Уникальное название
      if(!this.editing && this.strGames.some(c => c.name === game.name)) {
        errors.push("Такое название уже существует")
      }

      // Проверка букв и пробелов
      const regex = /^[a-zA-Z0-9\s]+$/;
      if(game.name && !regex.test(game.name)) errors.push("Поле 'Название' может состоять только из букв и пробелов");
      if(game.country && !regex.test(game.studio)) errors.push("Поле 'Студия' может состоять только из букв и пробелов")

      this.errorMessage = errors;

      return errors.length === 0;
    },

    async fetchGames() {
      const res = await api.get("/strGames");
      this.strGames = res.data;
    },
    async saveGame() {
      const payload = {
        ...this.newGame,
        year: Number(this.newGame.year),
        sold: Number(this.newGame.sold),
      };
      
      // Валидация
      if(!this.validateGame(payload)) return;

      try {
        if(this.editing) {
          // Редакутирование
          await api.put(`/strGames/${this.editing.id}`, payload);
          this.editing = null;
        } else {
          // Добавление
          await api.post("/strGames", payload);
        }

        // Очистка форм
        this.errorMessage = [];
        this.newGame = { name: "", studio: "", year: "", sold: "" };
        this.fetchGames();
      } catch(err) {
        this.errorMessage = ["Ошибка при сохранении"];
        console.log(err);
      }
    },
    editGame(game) {
      this.newGame = { ...game,}; 
      this.editing = game;
    },
    async deleteGame(id) {
      await api.delete(`/strGames/${id}`);
      this.fetchGames();
    },
  },
};
</script>

<style>
.container {
  padding: 20px;
  font-family: Arial, sans-serif;
}
.form {
  margin-bottom: 20px;
}
input {
  margin: 5px;
}
button {
  margin: 5px;
}
</style>
