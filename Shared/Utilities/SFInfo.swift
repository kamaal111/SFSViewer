//
//  SFInfo.swift
//  SFSViewer
//
//  Created by Kamaal M Farah on 23/12/2021.
//

import Foundation

public struct SFInfo {
    private let names: [SFName]
    private let supportedVersions: [String: SFVersions]

    public init() throws {
        self.supportedVersions = try Self.decodeFile(withName: "supported_versions")
        self.names = try Self.decodeFile(withName: "names")
    }

    public enum InitializerErrors: Error {
        case resourceNotFound(name: String)
        case encodingError(error: Error)
        case decodingError(error: Error)
    }

    private static func decodeFile<T: Decodable>(withName name: String) throws -> T {
        guard let url = Bundle.main.url(forResource: name, withExtension: "json")
        else { throw InitializerErrors.resourceNotFound(name: name) }

        let data: Data
        do {
            data = try Data(contentsOf: url)
        } catch {
            throw InitializerErrors.encodingError(error: error)
        }

        let decodedData: T
        do {
            decodedData = try JSONDecoder().decode(T.self, from: data)
        } catch {
            throw InitializerErrors.decodingError(error: error)
        }

        return decodedData
    }
}

struct SFName: Codable, Hashable {
    let name: String
    let releaseYear: String

    init(name: String, releaseYear: String) {
        self.name = name
        self.releaseYear = releaseYear
    }

    enum CodingKeys: String, CodingKey {
        case name
        case releaseYear = "release_year"
    }
}

struct SFVersions: Codable, Hashable {
    let iOS: String
    let macOS: String
    let tvOS: String
    let watchOS: String
}
