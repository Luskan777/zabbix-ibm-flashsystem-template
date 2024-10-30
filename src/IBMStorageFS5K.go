package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type IBMStorageFS5K struct {
	hostname string
	port     string
	user     string
	password string
}

type IBMStorage interface {
	discoverEnclosures()
	discoverDrives()
	discoverArrays()
	discoverVolumes()
}

func (s *IBMStorageFS5K) handleCommand(cmd CommandStruct) {
	if cmd.isDescoveryCmd {
		// Commands that are used to discover the storage
		switch cmd.command {
		case "discoverenclosures":
			s.discoverEnclosures()
		case "discovermdiskgroups":
			s.discoverMdiskGroup()
		case "discoverdrives":
			s.discoverDrives()
		case "discovermdisks":
			s.discoverArrays()
		case "discovervolumes":
			s.discoverVolumes()
		default:
			log.Fatalln("Invalid command, please use discoverenclosures, discovermdiskgroups,  discoverdrives, discovermdisks or discovervdisks")
		}
	} else {
		switch cmd.commandArray[0] {
		case "getmdiskstatus":
			arrayID := cmd.commandArray[1]
			s.getArrayStatus(arrayID)
		case "getmdiskgroupstatus":
			mdiskGroupID := cmd.commandArray[1]
			s.getMdiskGroupStatus(mdiskGroupID)
		case "getmdiskgroupcapacity":
			mdiskGroupID := cmd.commandArray[1]
			s.getMdiskGroupCapacity(mdiskGroupID)
		case "getmdiskgroupfreecapacity":
			mdiskGroupID := cmd.commandArray[1]
			s.getMdiskGroupFreeCapacity(mdiskGroupID)
		case "getvolumestatus":
			volumeID := cmd.commandArray[1]
			s.getVolumeStatus(volumeID)
		case "getenclosurestatus":
			enclosureID := cmd.commandArray[1]
			s.getEnclosureStatus(enclosureID)
		case "getbatterystatus":
			enclosureID := cmd.commandArray[1]
			batteryID := cmd.commandArray[2]
			s.getBatteryStatus(enclosureID, batteryID)
		case "getnodestatus":
			enclosureID := cmd.commandArray[1]
			nodeID := cmd.commandArray[2]
			s.getNodeStatus(enclosureID, nodeID)
		case "getpsustatus":
			enclosureID := cmd.commandArray[1]
			psuID := cmd.commandArray[2]
			s.getPSUStatus(enclosureID, psuID)
		case "getdrivestatus":
			enclosureID := cmd.commandArray[1]
			driveID := cmd.commandArray[2]
			s.getDriveStatus(enclosureID, driveID)
		case "getdriveslotstatus":
			enclosureID := cmd.commandArray[1]
			driveID := cmd.commandArray[2]
			s.getDriveSlotStatus(enclosureID, driveID)
		default:
			log.Fatalln("Invalid command, please use getmdiskstatus, getmdiskgroupstatus, getmdiskgroupcapacity, getmdiskgroupfreecapacity, getvolumestatus, getenclosurestatus, getbatterystatus, getnodestatus, getpsustatus, getdrivestatus or getdriveslotstatus")
		}

	}
}

// Discover structs and methods

type discoverEnclosures struct {
	EnclosureID string `json:"{#ENCLOSURE_ID}"`
}

type discoverMdiskGroup struct {
	MdiskGroupID       string `json:"{#MDISK_GROUP_ID}"`
	MdiskGroupName     string `json:"{#MDISK_GROUP_NAME}"`
	MdiskGroupCapacity string `json:"{#MDISK_GROUP_CAPACITY}"`
	MdiskGroupEasyTier string `json:"{#MDISK_GROUP_EASY_TIER}"`
}

type discoverDrives struct {
	EnclosureID string `json:"{#ENCLOSURE_ID}"`
	DriveID     string `json:"{#DRIVE_ID}"`
}

type discoverArrays struct {
	ArrayID   string `json:"{#ARRAY_ID}"`
	ArrayName string `json:"{#ARRAY_NAME}"`
}

type discoverVolumes struct {
	VDISK_ID   string `json:"{#VDISK_ID}"`
	VDISK_NAME string `json:"{#VDISK_NAME}"`
}

func (s *IBMStorageFS5K) discoverEnclosures() {
	cmd := "lsenclosure -nohdr -delim :"
	result := executeCmd(s, cmd)
	resultSplitPerLine := parseResultPerline(result)

	enclosures := []discoverEnclosures{}

	for _, line := range resultSplitPerLine {
		if line == "" {
			continue

		}
		enclosure := discoverEnclosures{}
		enclosure.EnclosureID = formatResult(line, 0)
		enclosures = append(enclosures, enclosure)
	}

	enclosuresJSON, err := json.Marshal(enclosures)
	if err != nil {
		log.Fatalln("Failed to marshal enclosures: ", err)
	}

	fmt.Println(string(enclosuresJSON))
}

func (s *IBMStorageFS5K) discoverMdiskGroup() {
	cmd := "lsmdiskgrp -nohdr -delim :"
	result := executeCmd(s, cmd)
	resultSplitPerLine := parseResultPerline(result)

	mdiskGroups := []discoverMdiskGroup{}

	for _, line := range resultSplitPerLine {
		if line == "" {
			continue
		}
		mdiskGroup := discoverMdiskGroup{}
		mdiskGroup.MdiskGroupID = formatResult(line, 0)
		mdiskGroup.MdiskGroupName = formatResult(line, 1)
		mdiskGroup.MdiskGroupCapacity = formatResult(line, 5)
		mdiskGroup.MdiskGroupEasyTier = formatResult(line, 14)
		mdiskGroups = append(mdiskGroups, mdiskGroup)
	}

	mdiskGroupsJSON, err := json.Marshal(mdiskGroups)
	if err != nil {
		log.Fatalln("Failed to marshal mdiskGroups: ", err)
	}

	fmt.Println(string(mdiskGroupsJSON))
}

func (s *IBMStorageFS5K) discoverDrives() {
	cmd := "lsdrive -nohdr -delim :"
	result := executeCmd(s, cmd)
	resultSplitPerLine := parseResultPerline(result)

	drives := []discoverDrives{}

	for _, line := range resultSplitPerLine {
		if line == "" {
			continue
		}
		driveStatus := formatResult(line, 1)
		if driveStatus != "online" {
			continue
		}

		drive := discoverDrives{}

		drive.EnclosureID = formatResult(line, 9)
		drive.DriveID = formatResult(line, 0)
		drives = append(drives, drive)
	}

	drivesJSON, err := json.Marshal(drives)
	if err != nil {
		log.Fatalln("Failed to marshal drives: ", err)
	}

	fmt.Println(string(drivesJSON))
}

func (s *IBMStorageFS5K) discoverArrays() {
	cmd := "lsarray -nohdr -delim :"
	result := executeCmd(s, cmd)
	resultSplitPerLine := parseResultPerline(result)

	arrays := []discoverArrays{}

	for _, line := range resultSplitPerLine {
		if line == "" {
			continue
		}
		array := discoverArrays{}
		array.ArrayID = formatResult(line, 0)
		array.ArrayName = formatResult(line, 1)
		arrays = append(arrays, array)
	}

	arraysJSON, err := json.Marshal(arrays)
	if err != nil {
		log.Fatalln("Failed to marshal arrays: ", err)
	}

	fmt.Println(string(arraysJSON))
}

func (s *IBMStorageFS5K) discoverVolumes() {
	cmd := "lsvdisk -nohdr -delim :"
	result := executeCmd(s, cmd)
	resultSplitPerLine := parseResultPerline(result)

	volumes := []discoverVolumes{}

	for _, line := range resultSplitPerLine {
		if line == "" {
			continue
		}
		volume := discoverVolumes{}
		volume.VDISK_ID = formatResult(line, 0)
		volume.VDISK_NAME = formatResult(line, 1)
		volumes = append(volumes, volume)
	}

	volumesJSON, err := json.Marshal(volumes)
	if err != nil {
		log.Fatalln("Failed to marshal volumes: ", err)
	}

	fmt.Println(string(volumesJSON))
}

// Items states structs and methods

// state of arrays
func (s *IBMStorageFS5K) getArrayStatus(arrayID string) {
	cmd := fmt.Sprintf("lsarray -filtervalue mdisk_id=%s -delim : -nohdr", arrayID)
	result := executeCmd(s, cmd)

	if result == "" {
		log.Fatalln("Array not found")
	}

	arrayStatus := formatResult(result, 2)

	if arrayStatus == "online" {
		arrayStatus = "1"
	} else {
		arrayStatus = "0"
	}

	fmt.Println(arrayStatus)
}

// State of mdisk groups
func (s *IBMStorageFS5K) getMdiskGroupStatus(mdiskGroupID string) {
	cmd := fmt.Sprintf("lsmdiskgrp -filtervalue id=%s -delim : -nohdr", mdiskGroupID)
	result := executeCmd(s, cmd)

	if result == "" {
		log.Fatalln("MdiskGroup not found")
	}

	mdiskGroupStatus := formatResult(result, 2)

	if mdiskGroupStatus == "online" {
		mdiskGroupStatus = "1"
	} else {
		mdiskGroupStatus = "0"
	}

	fmt.Println(mdiskGroupStatus)
}

// mdiskgroup capacity
func (s *IBMStorageFS5K) getMdiskGroupCapacity(mdiskGroupID string) {
	cmd := fmt.Sprintf("lsmdiskgrp -filtervalue id=%s -delim : -nohdr", mdiskGroupID)
	result := executeCmd(s, cmd)

	if result == "" {
		log.Fatalln("MdiskGroup not found")
	}

	mdiskGroupCapacity := formatResult(result, 5)

	mdiskGroupsize := getValueDataSizeWithoutUnit(mdiskGroupCapacity)

	fmt.Println(mdiskGroupsize)
}

// mdiskgroup free capacity
func (s *IBMStorageFS5K) getMdiskGroupFreeCapacity(mdiskGroupID string) {
	cmd := fmt.Sprintf("lsmdiskgrp -filtervalue id=%s -delim : -nohdr", mdiskGroupID)
	result := executeCmd(s, cmd)

	if result == "" {
		log.Fatalln("MdiskGroup not found")
	}

	mdiskGroupFreeCapacity := formatResult(result, 7)

	mdiskGroupFreeSize := getValueDataSizeWithoutUnit(mdiskGroupFreeCapacity)

	fmt.Println(mdiskGroupFreeSize)
}

// state of volumes
func (s *IBMStorageFS5K) getVolumeStatus(volumeID string) {
	cmd := fmt.Sprintf("lsvdisk -filtervalue id=%s -delim : -nohdr", volumeID)
	result := executeCmd(s, cmd)

	if result == "" {
		log.Fatalln("Volume not found")
	}

	volumeStatus := formatResult(result, 4)

	if volumeStatus == "online" {
		volumeStatus = "1"
	} else {
		volumeStatus = "0"
	}

	fmt.Println(volumeStatus)
}

// state of enclosure
func (s *IBMStorageFS5K) getEnclosureStatus(enclosureID string) {
	cmd := fmt.Sprintf("lsenclosure -filtervalue id=%s -delim : -nohdr", enclosureID)
	result := executeCmd(s, cmd)

	if result == "" {
		log.Fatalln("Enclosure not found")
	}

	enclosureStatus := formatResult(result, 1)

	if enclosureStatus == "online" {
		enclosureStatus = "1"
	} else {
		enclosureStatus = "0"
	}

	fmt.Println(enclosureStatus)
}

// state of battery
func (s *IBMStorageFS5K) getBatteryStatus(enclosureID, batteryID string) {
	cmd := fmt.Sprintf("lsenclosurebattery -filtervalue enclosure_id=%s:battery_id=%s  -delim : -nohdr", enclosureID, batteryID)
	result := executeCmd(s, cmd)

	if result == "" {
		log.Fatalln("Battery not found")
	}

	batteryStatus := formatResult(result, 2)

	if batteryStatus == "online" {
		batteryStatus = "1"
	} else {
		batteryStatus = "0"
	}

	fmt.Println(batteryStatus)
}

// state of canister/nodes
func (s *IBMStorageFS5K) getNodeStatus(enclosureID, nodeID string) {
	cmd := fmt.Sprintf("lsenclosurecanister -filtervalue enclosure_id=%s:node_id=%s  -delim : -nohdr", enclosureID, nodeID)
	result := executeCmd(s, cmd)

	if result == "" {
		log.Fatalln("Node not found")
	}

	nodeStatus := formatResult(result, 2)

	if nodeStatus == "online" {
		nodeStatus = "1"
	} else {
		nodeStatus = "0"
	}

	fmt.Println(nodeStatus)
}

// state of psu's
func (s *IBMStorageFS5K) getPSUStatus(enclosureID, psuID string) {
	cmd := fmt.Sprintf("lsenclosurepsu -filtervalue enclosure_id=%s:PSU_id=%s  -delim : -nohdr", enclosureID, psuID)
	result := executeCmd(s, cmd)

	if result == "" {
		log.Fatalln("PSU not found")
	}

	psuStatus := formatResult(result, 2)

	if psuStatus == "online" {
		psuStatus = "1"
	} else {
		psuStatus = "0"
	}

	fmt.Println(psuStatus)
}

// state of drives
func (s *IBMStorageFS5K) getDriveStatus(enclosureID, driveID string) {
	cmd := fmt.Sprintf("lsdrive -filtervalue enclosure_id=%s:id=%s  -delim : -nohdr", enclosureID, driveID)
	result := executeCmd(s, cmd)

	if result == "" {
		log.Fatalln("Drive not found")
	}

	driveStatus := formatResult(result, 1)

	if driveStatus == "online" {
		driveStatus = "1"
	} else {
		driveStatus = "0"
	}

	fmt.Println(driveStatus)
}

// state of slot of drives (Only drives that are in a slot)
func (s *IBMStorageFS5K) getDriveSlotStatus(enclosureID, driveID string) {
	cmd := fmt.Sprintf("lsenclosureslot -filtervalue enclosure_id=%s:drive_id=%s  -delim : -nohdr", enclosureID, driveID)
	result := executeCmd(s, cmd)

	if result == "" {
		log.Fatalln("drive/slot not found")
	}

	drivePresetStatus := formatResult(result, 2)

	if drivePresetStatus == "online" {
		drivePresetStatus = "1"
	} else {
		drivePresetStatus = "0"
	}

	fmt.Println(drivePresetStatus)
}
