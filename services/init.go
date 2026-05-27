package services

import "project-MVP/db"

// DefaultStore'ы для обратной совместимости.
// Инициализируются функцией InitDefaultStores() после подключения к БД.
var (
	DefaultProjectStore   *ProjectStore
	DefaultTaskStore      *TaskStore
	DefaultTeamStore      *TeamStore
	DefaultRBACStore      *RBACStore
	DefaultCommentStore   *CommentStore
	DefaultSprintStore    *SprintStore
	DefaultTagStore       *TagStore
	DefaultDatasetStore           *DatasetStore
	DefaultExperimentStore        *ExperimentStore
	DefaultDatasetDependencyStore *DatasetDependencyStore
	DefaultUserStore              *UserStore
)

// InitDefaultStores создаёт глобальные инстансы Store'ов с актуальным
// подключением к БД. ДОЛЖНА вызываться после db.Connect().
func InitDefaultStores() {
	DefaultProjectStore = NewProjectStore(db.DB)
	DefaultTaskStore = NewTaskStore(db.DB)
	DefaultTeamStore = NewTeamStore(db.DB)
	DefaultRBACStore = NewRBACStore(db.DB)
	DefaultCommentStore = NewCommentStore(db.DB)
	DefaultSprintStore = NewSprintStore(db.DB)
	DefaultTagStore = NewTagStore(db.DB)
	DefaultDatasetStore = NewDatasetStore(db.DB)
	DefaultExperimentStore = NewExperimentStore(db.DB)
	DefaultDatasetDependencyStore = NewDatasetDependencyStore(db.DB)
	DefaultUserStore = NewUserStore(db.DB)
}
