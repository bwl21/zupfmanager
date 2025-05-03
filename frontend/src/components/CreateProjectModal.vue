<template>
  <div class="modal" v-if="show">
    <div class="modal-content">
      <span class="close" @click="close">&times;</span>
      <h2>Neues Projekt erstellen</h2>
      <div class="form-group">
        <label for="title">Titel:</label>
        <input type="text" id="title" v-model="title" class="form-control">
      </div>
      <button @click="createProject" class="btn btn-primary">Erstellen</button>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    show: Boolean,
  },
  data() {
    return {
      title: '',
    };
  },
  methods: {
    createProject() {
      if (this.title.trim() !== '') {
        this.$emit('projectCreated', this.title);
        this.title = '';
        this.close();
      }
    },
    close() {
      this.$emit('close');
    },
  },
};
</script>

<style scoped>
.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal-content {
  background-color: white;
  padding: 20px;
  border-radius: 5px;
  width: 500px;
  position: relative;
}

.close {
  position: absolute;
  top: 10px;
  right: 10px;
  font-size: 20px;
  cursor: pointer;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
}

.form-control {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.btn {
  padding: 0.375rem 0.75rem;
  border-radius: 0.25rem;
  border: none;
  color: white;
  cursor: pointer;
}

.btn-primary {
  background-color: #007bff;
}
</style>
