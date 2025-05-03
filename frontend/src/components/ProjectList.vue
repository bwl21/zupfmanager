<template>
  <h1>Projektliste</h1>

  <button @click="openCreateProjectModal" class="btn btn-primary">Neues Projekt erstellen</button>

  <table class="table">
    <thead>
      <tr>
        <th>Titel</th>
        <th>Aktionen</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="project in projects" :key="project.project_id">
        <td>{{ project.project_title }}</td>
        <td>
          <button @click="editProject(project)" class="btn btn-secondary btn-sm">Bearbeiten</button>
          <button @click="deleteProject(project)" class="btn btn-danger btn-sm">Löschen</button>
          <button @click="viewProject(project)" class="btn btn-info btn-sm">Anzeigen</button>
        </td>
      </tr>
    </tbody>
  </table>

  <CreateProjectModal :show="showCreateProjectModal" @close="closeCreateProjectModal" @projectCreated="fetchProjects" />
</template>

<script>
import CreateProjectModal from './CreateProjectModal.vue';

export default {
  components: {
    CreateProjectModal,
  },
  data() {
    return {
      projects: [],
      showCreateProjectModal: false,
    };
  },
  async mounted() {
    await this.fetchProjects();
  },
  methods: {
    async fetchProjects() {
      this.projects = await window.go.main.App.GetAllProjects();
    },
    openCreateProjectModal() {
      this.showCreateProjectModal = true;
    },
    closeCreateProjectModal() {
      this.showCreateProjectModal = false;
    },
    async createProject(title) {
      await window.go.main.App.CreateProject(title);
      await this.fetchProjects();
    },
    editProject(project) {
      // TODO: Implement edit functionality
      console.log('Edit project', project);
    },
    async deleteProject(project) {
      // TODO: Implement delete functionality
	  if (confirm(`Möchten Sie das Projekt "${project.project_title}" wirklich löschen?`)) {
        await window.go.main.App.DeleteProject(project.project_id);
        await this.fetchProjects();
      }
    },
    viewProject(project) {
      // TODO: Implement view functionality (navigate to project details)
      console.log('View project', project);
    },
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

.btn-secondary {
  background-color: #6c757d;
}

.btn-danger {
  background-color: #dc3545;
}

.btn-info {
  background-color: #17a2b8;
}

.btn-sm {
  padding: 0.25rem 0.5rem;
  font-size: 0.875rem;
  line-height: 1.5;
  border-radius: 0.2rem;
}
</style>
