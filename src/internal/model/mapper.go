package model

import "github.com/jinzhu/copier"

func ToDomain[ModelType, DomainType any](modelObj *ModelType) (*DomainType, error) {
	var result DomainType
	if err := copier.Copy(&result, modelObj); err != nil {
		return nil, err
	}
	return &result, nil
}

func ToModel[ModelType, DomainType any](domainObj *DomainType) (*ModelType, error) {
	var result ModelType
	if err := copier.Copy(&result, domainObj); err != nil {
		return nil, err
	}
	return &result, nil
}

func ToDomainSlice[ModelType, DomainType any](modelObjs []ModelType) ([]*DomainType, error) {
	result := make([]*DomainType, len(modelObjs))
	var err error

	for i, modelObj := range modelObjs {
		result[i], err = ToDomain[ModelType, DomainType](&modelObj)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
