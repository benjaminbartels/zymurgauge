<template>
  <v-container>
       <v-btn color="login" @click="login()">Login</v-btn>
  </v-container>
</template>

<script>

export default {
  name: 'auth',
  props: ['auth'],
  methods: {
    login: function () {
      this.auth.login()
    }
  }
}
</script>
