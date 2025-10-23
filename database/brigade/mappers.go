package brigade

import (
	"brigade-service/service/brigade"
)

func MapInspectorsToMembers(inspectors []brigade.Inspector, brigadeID int) []Member {
	result := make([]Member, 0, len(inspectors))
	for _, i := range inspectors {
		result = append(result, MapInspectorToMember(i, brigadeID))
	}

	return result
}

func MapInspectorToMember(i brigade.Inspector, brigadeID int) Member {
	return Member{
		BrigadeID:   brigadeID,
		InspectorID: i.ID,
	}
}

func MapBrigadesFromDB(brigades []Brigade) []brigade.Brigade {
	result := make([]brigade.Brigade, 0, len(brigades))
	for _, b := range brigades {
		result = append(result, MapBrigadeFromDB(b))
	}

	return result
}

func MapBrigadeFromDB(b Brigade) brigade.Brigade {
	return brigade.Brigade{
		ID:        b.ID,
		Status:    brigade.Status(b.Status),
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}

func MapMembersToInspectors(members []Member) []brigade.Inspector {
	result := make([]brigade.Inspector, 0, len(members))
	for _, m := range members {
		result = append(result, MapMemberToInspector(m))
	}

	return result
}

func MapMemberToInspector(m Member) brigade.Inspector {
	return brigade.Inspector{
		ID:         m.InspectorID,
		AssignedAt: m.AssignedAt,
	}
}
