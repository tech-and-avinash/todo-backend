package repositories

type FileRepository struct{}

func NewFileRepository() *FileRepository {
	return &FileRepository{}
}
