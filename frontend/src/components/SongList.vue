<template>
  <h1>Songliste</h1>
  <table class="table">
    <thead>
      <tr>
        <th>Titel</th>
        <th>Dateiname</th>
        <th>Genre</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="song in songs" :key="song.id">
        <td>{{ song.title }}</td>
        <td>{{ song.filename }}</td>
        <td>{{ song.genre }}</td>
      </tr>
    </tbody>
  </table>
</template>

<script>
export default {
  data() {
    return {
      songs: [],
    };
  },
  async mounted() {
    this.songs = await window.go.main.App.GetAllSongs();
  },
};
</script>

<style scoped>
.table {
  width: 100%;
  border-collapse: collapse;
}

.table th,
.table td {
  border: 1px solid #ddd;
  padding: 8px;
  text-align: left;
}

.table th {
  background-color: #f2f2f2;
}
</style>
