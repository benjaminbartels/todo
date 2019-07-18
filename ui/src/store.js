import Vue from 'vue'
import Vuex from 'vuex'
import { HTTP } from "./http-common";

Vue.use(Vuex)

const state = {
  todos: '[]',
  errorMsg: ''
}

const mutations = {
  loadTodos(state, todos) {
    state.todos = todos
  },
  addTodo(state, todo) {
    state.todos.push(todo)
  },
  removeTodo(state, todo) {
    state.todos.splice(state.todos.indexOf(todo), 1)
  },
  editTodo(state, { todo, title = todo.title, completed = todo.completed }) {
    todo.title = title
    todo.completed = completed
  },
  populateError(state, errorMsg) {
    state.errorMsg = errorMsg
  }
}

const actions = {
  loadTodos({ commit }) {
    HTTP
      .get('/todos')
      .then(r => {
        commit('loadTodos', r.data)
      })
      .catch(e => {
        populateError(commit,e)
      })
  },
  addTodo({ commit }, title) {
    HTTP
      .post('/todos', { title })
      .then(r => {
        commit('addTodo', r.data)
      })
      .catch(e => {
        populateError(commit,e)
      })
  },
  removeTodo({ commit }, todo) {
    HTTP
      .delete('/todos/' + todo.id)
      .then(
        commit('removeTodo', todo)
      )
      .catch(e => {
        populateError(commit,e)
      })
  },
  toggleTodo({ commit }, todo) {
    HTTP
      .put('/todos/' + todo.id, todo)
      .then(
        commit('editTodo', { todo, completed: !todo.completed })
      )
      .catch(e => {
        populateError(commit,e)
      })
  },
  editTodo({ commit }, { todo, value }) {
    HTTP
      .put('/todos/' + todo.id, todo)
      .then(
        commit('editTodo', { todo, title: value })
      )
      .catch(e => {
        populateError(commit,e)
      })
  },
  toggleAll({ state, commit }, completed) {
    state.todos.forEach((todo) => {
      HTTP
        .put('/todos/' + todo.id, todo)
        .then(
          commit('editTodo', { todo, completed: completed })
        )
        .catch(e => {
          populateError(commit,e)
        })
    })
  },
  clearCompleted({ state, commit }) {
    state.todos.filter(todo => todo.completed)
      .forEach(todo => {
        HTTP
          .delete('/todos/' + todo.id)
          .then(
            commit('removeTodo', todo)
          )
          .catch(e => {
            populateError(commit,e)
          })
      })
  }
}

function populateError(commit, errorMsg) {
  commit('populateError', errorMsg)
  setTimeout(() => commit('populateError', ''), 5000)
}

export default new Vuex.Store({
  state,
  mutations,
  actions
})