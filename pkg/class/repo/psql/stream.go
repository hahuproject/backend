package class_psql_repo

import class_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/domain"

func (repo PsqlClassRepoAdapter) StoreStream(stream class_domain.Stream) (class_domain.Stream, error) {
	var addedStream class_domain.Stream

	err := repo.db.QueryRow(`INSERT INTO public.streams (name, department_id) VALUES ($1, $2) RETURNING stream_id`, stream.Name, stream.Department.ID).Scan(&stream.ID)
	if err != nil {
		return addedStream, err
	}

	return addedStream, nil
}

func (repo PsqlClassRepoAdapter) UpdateStream(stream class_domain.Stream) (class_domain.Stream, error) {
	var updatedStream class_domain.Stream
	err := repo.db.QueryRow("UPDATE public.streams SET name = $1 WHERE stream_id = $2 RETURNING * ", stream.Name, stream.ID).Scan(&updatedStream.ID, &updatedStream.Name, &updatedStream.Department.ID)
	return updatedStream, err
}

func (repo PsqlClassRepoAdapter) DeleteStream(streamId string) error {
	repo.log.Println(streamId)
	_, err := repo.db.Query("DELETE FROM public.streams WHERE stream_id = $1", streamId)
	return err
}
